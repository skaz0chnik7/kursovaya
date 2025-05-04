const imageInput = document.getElementById("image");
const uploadedImage = document.getElementById("uploadedImage");
const downloadLink = document.getElementById("downloadLink");
const infoBlock = document.getElementById("info");


let fileName = "";
let scale = 1; // Изначальный масштаб изображения


function updatePreview(blob) {
    const url = URL.createObjectURL(blob);
    uploadedImage.src = url;
    downloadLink.href = url;
    downloadLink.style.display = "block";
}

// Масштабирование
function scaleImage(delta) {
    scale += delta;
    if (scale < 0.1) scale = 0.1; // Минимальный масштаб
    uploadedImage.style.transform = `scale(${scale})`;
}

imageInput.addEventListener("change", function () {
    const file = this.files[0];
    if (!file || !file.type.startsWith("image/")) {
        alert("Пожалуйста, выберите изображение.");
        return;
    }

    const reader = new FileReader();
    reader.onload = (e) => {
        uploadedImage.src = e.target.result;
        document.getElementById("preview").style.display = "block";
    };
    reader.readAsDataURL(file);

    const formData = new FormData();
    formData.append("image", file);

    fetch("/upload", { method: "POST", body: formData })
        .then((res) => res.json())
        .then((data) => {
            fileName = data.fileName;
            console.log("File uploaded successfully:", fileName);
        })
        .catch((err) => console.error("Upload error:", err));
});

document.getElementById("resetButton").addEventListener("click", () => {
    if (!fileName) {
        alert("Сначала загрузите изображение.");
        return;
    }
    fetch(`/reset?file=${fileName}`)
        .then((res) => res.blob())
        .then(updatePreview)
        .catch((err) => console.error("Reset failed:", err));
});

document.getElementById("resizeButton").addEventListener("click", () => {
    if (!fileName) {
        alert("Сначала загрузите изображение.");
        return;
    }

    const width = document.getElementById("width").value;
    const height = document.getElementById("height").value;

    if (!width || !height) {
        alert("Укажите ширину и высоту для изменения размера.");
        return;
    }

    fetch(`/resize?file=${fileName}&width=${width}&height=${height}`)
        .then((res) => res.blob())
        .then(updatePreview)
        .catch((err) => console.error("Resize failed:", err));
});

document.getElementById("rotateButton").addEventListener("click", () => {
    if (!fileName) {
        alert("Сначала загрузите изображение.");
        return;
    }

    const rotate = document.getElementById("rotate").value;

    if (!rotate) {
        alert("Введите угол поворота.");
        return;
    }

    fetch(`/rotate?file=${fileName}&rotate=${rotate}`)
        .then((res) => res.blob())
        .then(updatePreview)
        .catch((err) => console.error("Rotate failed:", err));
});

document.getElementById("cropButton").addEventListener("click", () => {
    if (!fileName) {
        alert("Сначала загрузите изображение.");
        return;
    }

    const cropX = document.getElementById("cropX").value;
    const cropY = document.getElementById("cropY").value;
    const cropWidth = document.getElementById("cropWidth").value;
    const cropHeight = document.getElementById("cropHeight").value;

    if (!cropX || !cropY || !cropWidth || !cropHeight) {
        alert("Заполните все поля для обрезки.");
        return;
    }

    fetch(`/crop?file=${fileName}&x=${cropX}&y=${cropY}&width=${cropWidth}&height=${cropHeight}`)
        .then((res) => res.blob())
        .then(updatePreview)
        .catch((err) => console.error("Crop failed:", err));
});

document.getElementById("filterButton").addEventListener("click", () => {
    if (!fileName) {
        alert("Сначала загрузите изображение.");
        return;
    }

    const filter = document.getElementById("filterMode").value;

    fetch(`/filter?file=${fileName}&filter=${filter}`)
        .then((res) => res.blob())
        .then(updatePreview)
        .catch((err) => console.error("Filter failed:", err));
});

document.getElementById("convertButton").addEventListener("click", () => {
    if (!fileName) {
        alert("Сначала загрузите изображение.");
        return;
    }

    const format = document.getElementById("convertFormat").value;

    fetch(`/convert?file=${fileName}&format=${format}`)
        .then((res) => res.blob())
        .then((blob) => {
            updatePreview(blob);
            const newFileName = `${fileName.split('.')[0]}.${format}`;
            downloadLink.download = newFileName;
            console.log("Converted file name:", newFileName);
        })
        .catch((err) => console.error("Convert failed:", err));
});

// Информация об изображении
document.getElementById("infoButton").addEventListener("click", () => {
    if (!fileName) {
        alert("Сначала загрузите изображение.");
        return;
    }

    fetch(`/info?file=${fileName}`)
        .then((res) => {
            if (!res.ok) {
                throw new Error(`Ошибка: ${res.status} ${res.statusText}`);
            }
            return res.json();
        })
        .then((info) => {
            infoBlock.innerHTML = `
                <p>Имя файла: ${info.file_name}</p>
                <p>Размер: ${(info.size_bytes / 1024).toFixed(2)} KB</p>
                <p>Дата последнего изменения: ${new Date(info.last_modified).toLocaleString()}</p>
                <p>Ширина: ${info.width}px</p>
                <p>Высота: ${info.height}px</p>
                <p>Формат: ${info.format}</p>`;
        })
        .catch((err) => {
            console.error("Info fetch failed:", err);
            alert("Не удалось получить информацию об изображении.");
        });
});

