  document.getElementById("form").addEventListener("submit", async (e) => {
      e.preventDefault()
      
      const data = new FormData(e.target)
      
      const res = await fetch("/upload",{
        method: "POST",
        body: data
      })
      if (!res.ok) {
          const err = await res.json();
          document.getElementById("link").textContent = err.error;
          return;
      }
      const str = await res.text()
      
      document.getElementById("link").textContent = `File link: <a href="${str}">${str}</a>`
})