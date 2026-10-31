const createForm = document.querySelector("#create-form");
const lookupForm = document.querySelector("#lookup-form");
const urlInput = document.querySelector("#url");
const lookupInput = document.querySelector("#lookup-id");
const createButton = document.querySelector("#create-button");
const emptyState = document.querySelector("#empty-state");
const qrImage = document.querySelector("#qr-image");
const loader = document.querySelector("#loader");
const resultTitle = document.querySelector("#result-title");
const resultBadge = document.querySelector("#result-badge");
const resultDetails = document.querySelector("#result-details");
const qrId = document.querySelector("#qr-id");
const notice = document.querySelector("#notice");
const copyButton = document.querySelector("#copy-button");
const downloadButton = document.querySelector("#download-button");

let currentImageUrl = null;
let currentId = null;

function showNotice(message, success = false) {
  notice.textContent = message;
  notice.className = success ? "notice success" : "notice";
  notice.hidden = false;
}

function clearNotice() {
  notice.hidden = true;
  notice.textContent = "";
}

function setLoading(loading) {
  createButton.disabled = loading;
  loader.hidden = !loading;
  if (loading) {
    emptyState.hidden = true;
    qrImage.hidden = true;
    resultDetails.hidden = true;
    resultTitle.textContent = "Обрабатываем запрос…";
    resultBadge.textContent = "ЗАГРУЗКА";
  }
}

function displayImage(url, id, title = "QR-код готов") {
  if (currentImageUrl && currentImageUrl.startsWith("blob:")) {
    URL.revokeObjectURL(currentImageUrl);
  }

  currentImageUrl = url;
  currentId = id;
  qrImage.src = url;
  qrImage.hidden = false;
  emptyState.hidden = true;
  resultDetails.hidden = false;
  qrId.textContent = id;
  resultTitle.textContent = title;
  resultBadge.textContent = "ГОТОВО";
  clearNotice();
}

createForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  clearNotice();

  const url = urlInput.value.trim();
  if (!url) {
    showNotice("Введи ссылку для создания QR-кода.");
    return;
  }

  setLoading(true);

  try {
    const response = await fetch("/gen", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url }),
    });

    const data = await response.json().catch(() => ({}));
    if (!response.ok) {
      throw new Error(data.error || `Ошибка сервера: ${response.status}`);
    }

    if (!data.id || !data.image) {
      throw new Error("Сервер вернул неполный ответ: отсутствует id или image.");
    }

    displayImage(`data:image/png;base64,${data.image}`, data.id);
    lookupInput.value = data.id;
    showNotice("QR-код создан и сохранён в Redis.", true);
  } catch (error) {
    resultTitle.textContent = "Не удалось создать QR";
    resultBadge.textContent = "ОШИБКА";
    emptyState.hidden = false;
    showNotice(error.message || "Неизвестная ошибка.");
  } finally {
    setLoading(false);
  }
});

lookupForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  clearNotice();

  const id = lookupInput.value.trim();
  if (!id) {
    showNotice("Вставь UUID сохранённого QR-кода.");
    return;
  }

  resultTitle.textContent = "Загружаем QR-код…";
  resultBadge.textContent = "ЗАГРУЗКА";
  emptyState.hidden = true;
  qrImage.hidden = true;
  resultDetails.hidden = true;
  loader.hidden = false;

  try {
    const response = await fetch(`/code/${encodeURIComponent(id)}`);
    if (!response.ok) {
      const data = await response.json().catch(() => ({}));
      throw new Error(data.error || `QR-код не найден (${response.status}).`);
    }

    const blob = await response.blob();
    if (blob.type !== "image/png") {
      throw new Error("Сервер вернул данные, которые не являются PNG.");
    }

    displayImage(URL.createObjectURL(blob), id, "QR-код загружен из Redis");
    showNotice("Изображение успешно получено по UUID.", true);
  } catch (error) {
    resultTitle.textContent = "Не удалось получить QR";
    resultBadge.textContent = "ОШИБКА";
    emptyState.hidden = false;
    showNotice(error.message || "Неизвестная ошибка.");
  } finally {
    loader.hidden = true;
  }
});

copyButton.addEventListener("click", async () => {
  if (!currentId) return;
  try {
    await navigator.clipboard.writeText(currentId);
    showNotice("UUID скопирован.", true);
  } catch {
    showNotice("Браузер не разрешил доступ к буферу обмена.");
  }
});

downloadButton.addEventListener("click", () => {
  if (!currentImageUrl) return;
  const link = document.createElement("a");
  link.href = currentImageUrl;
  link.download = `qr-${currentId || "code"}.png`;
  document.body.appendChild(link);
  link.click();
  link.remove();
});


async function checkServiceHealth() {
  const status = document.querySelector("#service-status");
  try {
    const response = await fetch("/health", { cache: "no-store" });
    if (!response.ok) throw new Error("health check failed");
    const data = await response.json();
    status.innerHTML = `<span class="status-dot"></span> ${data.status === "ok" ? "Сервис работает" : "Сервис отвечает"}`;
  } catch {
    status.innerHTML = '<span class="status-dot offline"></span> Нет связи с API';
  }
}

checkServiceHealth();
