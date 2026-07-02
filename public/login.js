const form = document.querySelector("#loginForm");
const statusEl = document.querySelector("#loginStatus");

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  statusEl.textContent = "Autenticando...";

  const payload = {
    usuario: document.querySelector("#usuario").value.trim(),
    senha: document.querySelector("#senha").value,
  };

  try {
    const response = await fetch("/api/v1/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    const data = await response.json();

    if (!response.ok) {
      statusEl.textContent = data.erro || "Credenciais invalidas.";
      return;
    }

    localStorage.setItem("access_token", data.access_token);
    localStorage.setItem("refresh_token", data.refresh_token);
    localStorage.setItem("token_type", data.token_type);
    window.location.href = "/";
  } catch (error) {
    statusEl.textContent = "Nao foi possivel conectar com a API.";
  }
});
