const state = {
  alunos: [],
  eventos: [],
};

const publicStudentsPath = "/api/v1/public/alunos";

const els = {
  sessionStatus: document.querySelector("#sessionStatus"),
  loginLink: document.querySelector("#loginLink"),
  logoutButton: document.querySelector("#logoutButton"),
  adminPanel: document.querySelector("#admin"),
  studentCount: document.querySelector("#studentCount"),
  eventCount: document.querySelector("#eventCount"),
  studentsGrid: document.querySelector("#studentsGrid"),
  studentsStatus: document.querySelector("#studentsStatus"),
  eventsTimeline: document.querySelector("#eventsTimeline"),
  eventsStatus: document.querySelector("#eventsStatus"),
  adminStatus: document.querySelector("#adminStatus"),
  studentDialog: document.querySelector("#studentDialog"),
  dialogPhoto: document.querySelector("#dialogPhoto"),
  dialogName: document.querySelector("#dialogName"),
  dialogClass: document.querySelector("#dialogClass"),
  dialogPhotos: document.querySelector("#dialogPhotos"),
  toastMessage: document.querySelector("#toastMessage"),
};

let toastTimer;

document.addEventListener("DOMContentLoaded", () => {
  bindSession();
  bindForms();
  bindRefreshButtons();
  updateSessionUI();
  loadStudents();
  loadEvents();
});

function bindSession() {
  els.logoutButton.addEventListener("click", () => {
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("token_type");
    updateSessionUI();
  });

  document.querySelector("#closeStudentDialog").addEventListener("click", () => {
    els.studentDialog.close();
  });
}

function bindRefreshButtons() {
  document.querySelector("#refreshStudents").addEventListener("click", loadStudents);
  document.querySelector("#refreshEvents").addEventListener("click", loadEvents);
}

function bindForms() {
  document.querySelector("#studentForm").addEventListener("submit", async (event) => {
    event.preventDefault();
    const data = Object.fromEntries(new FormData(event.currentTarget).entries());
    const response = await apiFetch("/api/v1/alunos", {
      method: "POST",
      body: JSON.stringify(data),
    });
    await handleMutation(response, "Aluno criado com sucesso.");
    event.currentTarget.reset();
    loadStudents();
  });

  document.querySelector("#photoForm").addEventListener("submit", async (event) => {
    event.preventDefault();
    const formData = Object.fromEntries(new FormData(event.currentTarget).entries());
    const alunoID = formData.aluno_id;
    delete formData.aluno_id;

    const response = await apiFetch(`/api/v1/alunos/${alunoID}/fotos`, {
      method: "POST",
      body: JSON.stringify(formData),
    });
    await handleMutation(response, "Foto vinculada com sucesso.");
    event.currentTarget.reset();
    loadStudents();
  });

  document.querySelector("#eventForm").addEventListener("submit", async (event) => {
    event.preventDefault();
    const data = Object.fromEntries(new FormData(event.currentTarget).entries());
    const response = await apiFetch("/api/v1/eventos", {
      method: "POST",
      body: JSON.stringify(data),
    });
    const saved = await handleMutation(response, "Evento adicionado com sucesso.");
    if (saved) {
      event.currentTarget.reset();
      loadEvents();
    }
  });
}

async function loadStudents() {
  els.studentsStatus.textContent = "Carregando alunos...";
  try {
    const response = await fetch(publicStudentsPath);
    if (!response.ok) {
      throw new Error("Falha ao carregar alunos");
    }
    state.alunos = await response.json();
    els.studentCount.textContent = state.alunos.length;
    renderStudents();
    els.studentsStatus.textContent = state.alunos.length
      ? "Clique em um card para ver fotos aninhadas."
      : "Nenhum aluno cadastrado.";
  } catch (error) {
    els.studentsStatus.textContent = "Nao foi possivel carregar os alunos.";
  }
}

async function loadEvents() {
  els.eventsStatus.textContent = "Carregando eventos...";
  try {
    const response = await fetch("/api/v1/eventos");
    if (!response.ok) {
      throw new Error("Falha ao carregar eventos");
    }
    state.eventos = await response.json();
    els.eventCount.textContent = state.eventos.length;
    renderEvents();
    els.eventsStatus.textContent = state.eventos.length
      ? "Eventos ordenados pela API."
      : "Nenhum evento cadastrado.";
  } catch (error) {
    els.eventsStatus.textContent = "Nao foi possivel carregar os eventos.";
  }
}

function renderStudents() {
  els.studentsGrid.innerHTML = "";
  const isAdmin = hasToken();

  state.alunos.forEach((aluno) => {
    const card = document.createElement("article");
    card.className = "student-card";
    card.tabIndex = 0;

    const img = document.createElement("img");
    img.src = aluno.foto || "/assets/logos/logo-memorias-anuario-digital.png";
    img.alt = `Foto de ${aluno.nome}`;

    const name = document.createElement("h3");
    name.textContent = aluno.nome;

    const turma = document.createElement("p");
    turma.textContent = aluno.turma || "Turma nao informada";

    card.append(img, name, turma);

    if (isAdmin) {
      const actions = document.createElement("div");
      actions.className = "card-actions";

      const editButton = document.createElement("button");
      editButton.type = "button";
      editButton.textContent = "Editar";
      editButton.addEventListener("click", (event) => {
        event.stopPropagation();
        editStudent(aluno);
      });

      const deleteButton = document.createElement("button");
      deleteButton.type = "button";
      deleteButton.textContent = "Excluir";
      deleteButton.addEventListener("click", (event) => {
        event.stopPropagation();
        deleteStudent(aluno.id);
      });

      actions.append(editButton, deleteButton);
      card.append(actions);
    }

    card.addEventListener("click", () => openStudent(aluno.id));
    card.addEventListener("keydown", (event) => {
      if (event.key === "Enter") {
        openStudent(aluno.id);
      }
    });

    els.studentsGrid.append(card);
  });
}

function renderEvents() {
  els.eventsTimeline.innerHTML = "";

  state.eventos.forEach((evento) => {
    const card = document.createElement("article");
    card.className = "event-card";

    const time = document.createElement("time");
    time.dateTime = evento.data;
    time.textContent = formatDate(evento.data);

    const title = document.createElement("h3");
    title.textContent = evento.titulo;

    const description = document.createElement("p");
    description.textContent = evento.descricao;

    card.append(time, title, description);

    if (hasToken()) {
      const actions = document.createElement("div");
      actions.className = "card-actions";

      const editButton = document.createElement("button");
      editButton.type = "button";
      editButton.textContent = "Editar";
      editButton.addEventListener("click", () => editEvent(evento));

      const deleteButton = document.createElement("button");
      deleteButton.type = "button";
      deleteButton.textContent = "Excluir";
      deleteButton.addEventListener("click", () => deleteEvent(evento.id));

      actions.append(editButton, deleteButton);
      card.append(actions);
    }

    els.eventsTimeline.append(card);
  });
}

async function openStudent(id) {
  try {
    const response = await fetch(`${publicStudentsPath}/${id}`);
    if (!response.ok) {
      throw new Error("Aluno nao encontrado");
    }
    const aluno = await response.json();
    els.dialogPhoto.src = aluno.foto || "/assets/logos/logo-memorias-anuario-digital.png";
    els.dialogPhoto.alt = `Foto de ${aluno.nome}`;
    els.dialogName.textContent = aluno.nome;
    els.dialogClass.textContent = aluno.turma || "Turma nao informada";
    els.dialogPhotos.innerHTML = "";

    if (aluno.fotos && aluno.fotos.length) {
      aluno.fotos.forEach((foto) => {
        const img = document.createElement("img");
        img.src = foto.url;
        img.alt = foto.legenda || `Foto vinculada ao aluno ${aluno.nome}`;
        els.dialogPhotos.append(img);
      });
    } else {
      els.dialogPhotos.textContent = "Sem fotos vinculadas ainda.";
    }

    els.studentDialog.showModal();
  } catch (error) {
    els.studentsStatus.textContent = "Nao foi possivel abrir o perfil do aluno.";
  }
}

async function editStudent(aluno) {
  const nome = window.prompt("Novo nome do aluno:", aluno.nome);
  if (!nome) {
    return;
  }
  const turma = window.prompt("Turma:", aluno.turma || "2026.1") || aluno.turma;
  const foto = window.prompt("Foto URL:", aluno.foto || "") || aluno.foto;

  const response = await apiFetch(`/api/v1/alunos/${aluno.id}`, {
    method: "PUT",
    body: JSON.stringify({ nome, turma, foto }),
  });
  await handleMutation(response, "Aluno atualizado.");
  loadStudents();
}

async function deleteStudent(id) {
  const response = await apiFetch(`/api/v1/alunos/${id}`, { method: "DELETE" });
  await handleMutation(response, "Aluno removido.");
  loadStudents();
}

async function editEvent(evento) {
  const titulo = window.prompt("Novo titulo do evento:", evento.titulo);
  if (!titulo) {
    return;
  }
  const descricao = window.prompt("Nova descricao do evento:", evento.descricao);
  if (!descricao) {
    return;
  }
  const data = window.prompt("Data do evento (YYYY-MM-DD):", evento.data);
  if (!data) {
    return;
  }
  const imagemURL = window.prompt("Imagem URL:", evento.imagem_url || "") || evento.imagem_url || "";

  const response = await apiFetch(`/api/v1/eventos/${evento.id}`, {
    method: "PUT",
    body: JSON.stringify({
      titulo,
      descricao,
      data,
      imagem_url: imagemURL,
    }),
  });
  const saved = await handleMutation(response, "Evento atualizado com sucesso.");
  if (saved) {
    loadEvents();
  }
}

async function deleteEvent(id) {
  const response = await apiFetch(`/api/v1/eventos/${id}`, { method: "DELETE" });
  await handleMutation(response, "Evento removido.");
  loadEvents();
}

async function handleMutation(response, successMessage) {
  if (response.ok) {
    showMessage(successMessage);
    return true;
  }

  let data = {};
  try {
    data = await response.json();
  } catch (error) {
    data = {};
  }
  showMessage(data.erro || "A operacao falhou.", true);
  return false;
}

function showMessage(message, isError = false) {
  els.adminStatus.textContent = message;
  els.toastMessage.textContent = message;
  els.toastMessage.classList.toggle("toast-error", isError);
  els.toastMessage.classList.remove("hidden");

  window.clearTimeout(toastTimer);
  toastTimer = window.setTimeout(() => {
    els.toastMessage.classList.add("hidden");
    els.toastMessage.classList.remove("toast-error");
  }, 4200);
}

async function apiFetch(url, options = {}) {
  let response = await fetch(url, withAuth(options));
  if (response.status === 401 && (await refreshToken())) {
    response = await fetch(url, withAuth(options));
  }
  if (response.status === 401) {
    localStorage.removeItem("access_token");
    updateSessionUI();
  }
  return response;
}

function withAuth(options = {}) {
  const headers = new Headers(options.headers || {});
  headers.set("Content-Type", "application/json");
  const token = localStorage.getItem("access_token");
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }
  return { ...options, headers };
}

async function refreshToken() {
  const refresh = localStorage.getItem("refresh_token");
  if (!refresh) {
    return false;
  }

  const response = await fetch("/api/v1/auth/refresh", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token: refresh }),
  });
  if (!response.ok) {
    return false;
  }
  const data = await response.json();
  localStorage.setItem("access_token", data.access_token);
  localStorage.setItem("refresh_token", data.refresh_token);
  return true;
}

function updateSessionUI() {
  const authenticated = hasToken();
  els.sessionStatus.textContent = authenticated ? "Admin" : "Visitante";
  els.loginLink.classList.toggle("hidden", authenticated);
  els.logoutButton.classList.toggle("hidden", !authenticated);
  els.adminPanel.classList.toggle("hidden", !authenticated);
  renderStudents();
  renderEvents();
}

function hasToken() {
  return Boolean(localStorage.getItem("access_token"));
}

function formatDate(value) {
  if (!value) {
    return "Data nao informada";
  }
  const [year, month, day] = value.split("-");
  if (!year || !month || !day) {
    return value;
  }
  return `${day}/${month}/${year}`;
}
