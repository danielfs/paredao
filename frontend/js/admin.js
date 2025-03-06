// Base URL for API requests
const API_BASE_URL = 'http://localhost:8080';

// DOM Elements
let participantsTable;
let votacoesTable;
let participanteForm;
let votacaoForm;
let alertContainer;
let participanteModal;
let votacaoModal;
let votacaoParticipantesModal;
let votacaoParticipantesList;
let addParticipanteToVotacaoForm;

// Current state
let currentParticipanteId = null;
let currentVotacaoId = null;

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
  // Get DOM elements
  participantsTable = document.getElementById('participants-table');
  votacoesTable = document.getElementById('votacoes-table');
  participanteForm = document.getElementById('participante-form');
  votacaoForm = document.getElementById('votacao-form');
  alertContainer = document.getElementById('alert-container');
  participanteModal = document.getElementById('participante-modal');
  votacaoModal = document.getElementById('votacao-modal');
  votacaoParticipantesModal = document.getElementById('votacao-participantes-modal');
  votacaoParticipantesList = document.getElementById('votacao-participantes-list');
  addParticipanteToVotacaoForm = document.getElementById('add-participante-to-votacao-form');

  // Set up event listeners
  setupEventListeners();

  // Load initial data
  loadParticipantes();
  loadVotacoes();

  // Set up tabs
  setupTabs();
});

// Set up event listeners
function setupEventListeners() {
  // Participante form submission
  participanteForm.addEventListener('submit', (e) => {
    e.preventDefault();
    saveParticipante();
  });

  // Votacao form submission
  votacaoForm.addEventListener('submit', (e) => {
    e.preventDefault();
    saveVotacao();
  });

  // Add participante to votacao form submission
  if (addParticipanteToVotacaoForm) {
    addParticipanteToVotacaoForm.addEventListener('submit', (e) => {
      e.preventDefault();
      addParticipanteToVotacao();
    });
  }

  // Close buttons for modals
  document.querySelectorAll('.close').forEach(closeBtn => {
    closeBtn.addEventListener('click', () => {
      document.querySelectorAll('.modal').forEach(modal => {
        modal.style.display = 'none';
      });
    });
  });

  // Close modals when clicking outside
  window.addEventListener('click', (e) => {
    document.querySelectorAll('.modal').forEach(modal => {
      if (e.target === modal) {
        modal.style.display = 'none';
      }
    });
  });
}

// Set up tabs
function setupTabs() {
  const tabs = document.querySelectorAll('.tab');
  const tabContents = document.querySelectorAll('.tab-content');

  tabs.forEach(tab => {
    tab.addEventListener('click', () => {
      // Remove active class from all tabs and contents
      tabs.forEach(t => t.classList.remove('active'));
      tabContents.forEach(content => content.classList.remove('active'));

      // Add active class to clicked tab and corresponding content
      tab.classList.add('active');
      const contentId = tab.getAttribute('data-tab');
      document.getElementById(contentId).classList.add('active');
    });
  });
}

// Show alert message
function showAlert(message, type = 'success') {
  alertContainer.innerHTML = `
    <div class="alert alert-${type}">
      ${message}
    </div>
  `;

  // Auto-hide after 3 seconds
  setTimeout(() => {
    alertContainer.innerHTML = '';
  }, 3000);
}

// Load all participantes
async function loadParticipantes() {
  try {
    const response = await fetch(`${API_BASE_URL}/participantes`);
    if (!response.ok) {
      throw new Error('Failed to load participants');
    }

    const participantes = await response.json();
    renderParticipantesTable(participantes);
    populateParticipanteDropdown(participantes);
  } catch (error) {
    console.error('Error loading participants:', error);
    showAlert('Failed to load participants', 'danger');
  }
}

// Render participantes table
function renderParticipantesTable(participantes) {
  if (!participantsTable) return;

  let html = `
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Nome</th>
          <th>Foto</th>
          <th>Ações</th>
        </tr>
      </thead>
      <tbody>
  `;

  if (participantes.length === 0) {
    html += `
      <tr>
        <td colspan="4" style="text-align: center;">Nenhum participante encontrado</td>
      </tr>
    `;
  } else {
    participantes.forEach(participante => {
      html += `
        <tr>
          <td>${participante.id}</td>
          <td>${participante.nome}</td>
          <td>
            <img src="${participante.urlFoto}" alt="${participante.nome}" style="width: 50px; height: 50px; object-fit: cover; border-radius: 50%;">
          </td>
          <td>
            <button class="btn" onclick="editParticipante(${participante.id})">Editar</button>
            <button class="btn btn-danger" onclick="deleteParticipante(${participante.id})">Excluir</button>
          </td>
        </tr>
      `;
    });
  }

  html += `
      </tbody>
    </table>
    <button class="btn" onclick="openParticipanteModal()">Adicionar Participante</button>
  `;

  participantsTable.innerHTML = html;
}

// Populate participante dropdown for adding to votacao
function populateParticipanteDropdown(participantes) {
  const dropdown = document.getElementById('participante-id');
  if (!dropdown) return;

  dropdown.innerHTML = '<option value="">Selecione um participante</option>';
  
  participantes.forEach(participante => {
    dropdown.innerHTML += `<option value="${participante.id}">${participante.nome}</option>`;
  });
}

// Open participante modal for adding/editing
function openParticipanteModal(id = null) {
  currentParticipanteId = id;
  
  // Reset form
  participanteForm.reset();
  
  if (id) {
    // Edit mode - load participante data
    fetch(`${API_BASE_URL}/participantes/${id}`)
      .then(response => {
        if (!response.ok) throw new Error('Failed to load participante');
        return response.json();
      })
      .then(participante => {
        document.getElementById('participante-nome').value = participante.nome;
        document.getElementById('participante-url-foto').value = participante.urlFoto;
        document.getElementById('participante-modal-title').textContent = 'Editar Participante';
      })
      .catch(error => {
        console.error('Error loading participante:', error);
        showAlert('Failed to load participante', 'danger');
      });
  } else {
    // Add mode
    document.getElementById('participante-modal-title').textContent = 'Adicionar Participante';
  }
  
  participanteModal.style.display = 'block';
}

// Save participante (create or update)
function saveParticipante() {
  const nome = document.getElementById('participante-nome').value;
  const urlFoto = document.getElementById('participante-url-foto').value;
  
  if (!nome || !urlFoto) {
    showAlert('Nome e URL da foto são obrigatórios', 'danger');
    return;
  }
  
  const participante = {
    nome: nome,
    urlFoto: urlFoto
  };
  
  const url = currentParticipanteId 
    ? `${API_BASE_URL}/participantes/${currentParticipanteId}` 
    : `${API_BASE_URL}/participantes`;
  
  const method = currentParticipanteId ? 'PUT' : 'POST';
  
  fetch(url, {
    method,
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(participante)
  })
    .then(response => {
      if (!response.ok) throw new Error('Failed to save participante');
      return response.json();
    })
    .then(() => {
      participanteModal.style.display = 'none';
      showAlert(currentParticipanteId ? 'Participante atualizado com sucesso' : 'Participante criado com sucesso');
      loadParticipantes();
    })
    .catch(error => {
      console.error('Error saving participante:', error);
      showAlert('Failed to save participante', 'danger');
    });
}

// Delete participante
function deleteParticipante(id) {
  if (!confirm('Tem certeza que deseja excluir este participante?')) {
    return;
  }
  
  fetch(`${API_BASE_URL}/participantes/${id}`, {
    method: 'DELETE'
  })
    .then(response => {
      if (!response.ok) throw new Error('Failed to delete participante');
      showAlert('Participante excluído com sucesso');
      loadParticipantes();
    })
    .catch(error => {
      console.error('Error deleting participante:', error);
      showAlert('Failed to delete participante', 'danger');
    });
}

// Load all votacoes
async function loadVotacoes() {
  try {
    const response = await fetch(`${API_BASE_URL}/votacoes`);
    if (!response.ok) {
      throw new Error('Failed to load votacoes');
    }

    const votacoes = await response.json();
    renderVotacoesTable(votacoes);
  } catch (error) {
    console.error('Error loading votacoes:', error);
    showAlert('Failed to load votacoes', 'danger');
  }
}

// Render votacoes table
function renderVotacoesTable(votacoes) {
  if (!votacoesTable) return;

  let html = `
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Descrição</th>
          <th>Ações</th>
        </tr>
      </thead>
      <tbody>
  `;

  if (votacoes.length === 0) {
    html += `
      <tr>
        <td colspan="3" style="text-align: center;">Nenhuma votação encontrada</td>
      </tr>
    `;
  } else {
    votacoes.forEach(votacao => {
      html += `
        <tr>
          <td>${votacao.id}</td>
          <td>${votacao.descricao}</td>
          <td>
            <button class="btn" onclick="editVotacao(${votacao.id})">Editar</button>
            <button class="btn btn-danger" onclick="deleteVotacao(${votacao.id})">Excluir</button>
            <button class="btn" onclick="manageVotacaoParticipantes(${votacao.id})">Gerenciar Participantes</button>
          </td>
        </tr>
      `;
    });
  }

  html += `
      </tbody>
    </table>
    <button class="btn" onclick="openVotacaoModal()">Adicionar Votação</button>
  `;

  votacoesTable.innerHTML = html;
}

// Open votacao modal for adding/editing
function openVotacaoModal(id = null) {
  currentVotacaoId = id;
  
  // Reset form
  votacaoForm.reset();
  
  if (id) {
    // Edit mode - load votacao data
    fetch(`${API_BASE_URL}/votacoes/${id}`)
      .then(response => {
        if (!response.ok) throw new Error('Failed to load votacao');
        return response.json();
      })
      .then(votacao => {
        document.getElementById('votacao-descricao').value = votacao.descricao;
        document.getElementById('votacao-modal-title').textContent = 'Editar Votação';
      })
      .catch(error => {
        console.error('Error loading votacao:', error);
        showAlert('Failed to load votacao', 'danger');
      });
  } else {
    // Add mode
    document.getElementById('votacao-modal-title').textContent = 'Adicionar Votação';
  }
  
  votacaoModal.style.display = 'block';
}

// Save votacao (create or update)
function saveVotacao() {
  const descricao = document.getElementById('votacao-descricao').value;
  
  if (!descricao) {
    showAlert('Descrição é obrigatória', 'danger');
    return;
  }
  
  const votacao = {
    descricao: descricao
  };
  
  const url = currentVotacaoId 
    ? `${API_BASE_URL}/votacoes/${currentVotacaoId}` 
    : `${API_BASE_URL}/votacoes`;
  
  const method = currentVotacaoId ? 'PUT' : 'POST';
  
  fetch(url, {
    method,
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(votacao)
  })
    .then(response => {
      if (!response.ok) throw new Error('Failed to save votacao');
      return response.json();
    })
    .then(() => {
      votacaoModal.style.display = 'none';
      showAlert(currentVotacaoId ? 'Votação atualizada com sucesso' : 'Votação criada com sucesso');
      loadVotacoes();
    })
    .catch(error => {
      console.error('Error saving votacao:', error);
      showAlert('Failed to save votacao', 'danger');
    });
}

// Delete votacao
function deleteVotacao(id) {
  if (!confirm('Tem certeza que deseja excluir esta votação?')) {
    return;
  }
  
  fetch(`${API_BASE_URL}/votacoes/${id}`, {
    method: 'DELETE'
  })
    .then(response => {
      if (!response.ok) throw new Error('Failed to delete votacao');
      showAlert('Votação excluída com sucesso');
      loadVotacoes();
    })
    .catch(error => {
      console.error('Error deleting votacao:', error);
      showAlert('Failed to delete votacao', 'danger');
    });
}

// Manage votacao participantes
function manageVotacaoParticipantes(votacaoId) {
  currentVotacaoId = votacaoId;
  
  // Load votacao details
  fetch(`${API_BASE_URL}/votacoes/${votacaoId}`)
    .then(response => {
      if (!response.ok) throw new Error('Failed to load votacao');
      return response.json();
    })
    .then(votacao => {
      document.getElementById('votacao-participantes-title').textContent = `Participantes da Votação: ${votacao.descricao}`;
      
      // Load participantes for this votacao
      return fetch(`${API_BASE_URL}/votacoes/${votacaoId}/participantes`);
    })
    .then(response => {
      if (!response.ok) throw new Error('Failed to load participantes');
      return response.json();
    })
    .then(participantes => {
      renderVotacaoParticipantes(participantes);
      votacaoParticipantesModal.style.display = 'block';
    })
    .catch(error => {
      console.error('Error managing votacao participantes:', error);
      showAlert('Failed to manage votacao participantes', 'danger');
    });
}

// Render votacao participantes
function renderVotacaoParticipantes(participantes) {
  if (!votacaoParticipantesList) return;

  if (participantes.length === 0) {
    votacaoParticipantesList.innerHTML = '<p>Nenhum participante nesta votação</p>';
    return;
  }

  let html = '<ul class="participantes-list">';
  
  participantes.forEach(participante => {
    html += `
      <li>
        <div class="participant-card">
          <img src="${participante.urlFoto}" alt="${participante.nome}" class="participant-image">
          <div class="participant-info">
            <div class="participant-name">${participante.nome}</div>
          </div>
        </div>
      </li>
    `;
  });
  
  html += '</ul>';
  
  votacaoParticipantesList.innerHTML = html;
}

// Add participante to votacao
function addParticipanteToVotacao() {
  const participanteId = document.getElementById('participante-id').value;
  
  if (!participanteId) {
    showAlert('Selecione um participante', 'danger');
    return;
  }
  
  const data = {
    participanteId: parseInt(participanteId)
  };
  
  fetch(`${API_BASE_URL}/votacoes/${currentVotacaoId}/participantes`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })
    .then(response => {
      if (!response.ok) throw new Error('Failed to add participante to votacao');
      return response.json();
    })
    .then(() => {
      showAlert('Participante adicionado à votação com sucesso');
      document.getElementById('participante-id').value = '';
      
      // Refresh participantes list
      return fetch(`${API_BASE_URL}/votacoes/${currentVotacaoId}/participantes`);
    })
    .then(response => {
      if (!response.ok) throw new Error('Failed to load participantes');
      return response.json();
    })
    .then(participantes => {
      renderVotacaoParticipantes(participantes);
    })
    .catch(error => {
      console.error('Error adding participante to votacao:', error);
      showAlert('Failed to add participante to votacao', 'danger');
    });
}

// Make functions available globally
window.openParticipanteModal = openParticipanteModal;
window.editParticipante = openParticipanteModal;
window.deleteParticipante = deleteParticipante;
window.openVotacaoModal = openVotacaoModal;
window.editVotacao = openVotacaoModal;
window.deleteVotacao = deleteVotacao;
window.manageVotacaoParticipantes = manageVotacaoParticipantes;
