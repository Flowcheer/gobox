document.getElementById("form").addEventListener("submit", async (e) => {
  e.preventDefault();

  const data = new FormData(e.target);
  const xhr = new XMLHttpRequest();
  xhr.open("POST", "/upload", true);
  xhr.upload.onprogress = function (e) {
    if (e.lengthComputable) {
      const percent = Math.round((e.loaded / e.total) * 100);
      document.getElementById("link").textContent = `Uploading: ${percent}%`;
    }
  };
  xhr.onload = function () {
    if (xhr.status === 200) {
      document.getElementById("link").innerHTML =
        `File link: <a href="${xhr.responseText}">${xhr.responseText}</a>`;
    } else {
      document.getElementById("link").textContent = xhr.responseText;
    }
  };
  xhr.onerror = function () {
    document.getElementById("link").textContent = "Upload failed / Network error";
  };
  xhr.send(data);
});