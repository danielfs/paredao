// Base URL for API requests
const API_BASE_URL = 'http://localhost:8080';

// DOM Elements
let votacaoSelector;
let participantsContainer;
let voteButton;
let alertContainer;
let loadingIndicator;

// Current state
let currentVotacaoId = null;
let selectedParticipanteId = null;

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
  // Get DOM elements
  votacaoSelector = document.getElementById('votacao-selector');
  participantsContainer = document.getElementById('participants-container');
  voteButton = document.getElementById('vote-button');
  alertContainer = document.getElementById('alert-container');
  loadingIndicator = document.getElementById('loading-indicator');

  // Set up event listeners
  setupEventListeners();

  // Load initial data
  loadVotacoes();
});

// Set up event listeners
function setupEventListeners() {
  // Votacao selector change
  if (votacaoSelector) {
    votacaoSelector.addEventListener('change', () => {
      const votacaoId = votacaoSelector.value;
      if (votacaoId) {
        currentVotacaoId = parseInt(votacaoId);
        loadVotacaoParticipantes(currentVotacaoId);
      } else {
        currentVotacaoId = null;
        participantsContainer.innerHTML = '<p>Selecione uma votação para ver os participantes</p>';
        updateVoteButtonState();
      }
    });
  }

  // Vote button click
  if (voteButton) {
    voteButton.addEventListener('click', submitVote);
  }
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

// Show/hide loading indicator
function setLoading(isLoading) {
  if (loadingIndicator) {
    loadingIndicator.style.display = isLoading ? 'block' : 'none';
  }
}

// Load all votacoes
async function loadVotacoes() {
  if (!votacaoSelector) return;

  setLoading(true);
  
  try {
    const response = await fetch(`${API_BASE_URL}/votacoes`);
    if (!response.ok) {
      throw new Error('Failed to load votacoes');
    }

    const votacoes = await response.json();
    renderVotacaoSelector(votacoes);
  } catch (error) {
    console.error('Error loading votacoes:', error);
    showAlert('Falha ao carregar votações', 'danger');
  } finally {
    setLoading(false);
  }
}

// Render votacao selector
function renderVotacaoSelector(votacoes) {
  if (!votacaoSelector) return;

  votacaoSelector.innerHTML = '<option value="">Selecione uma votação</option>';
  
  if (votacoes.length === 0) {
    votacaoSelector.innerHTML += '<option disabled>Nenhuma votação disponível</option>';
    return;
  }
  
  votacoes.forEach(votacao => {
    console.log(votacao);
    votacaoSelector.innerHTML += `<option value="${votacao.id}">${votacao.descricao}</option>`;
  });
}

// Load participantes for a votacao
async function loadVotacaoParticipantes(votacaoId) {
  if (!participantsContainer) return;

  setLoading(true);
  selectedParticipanteId = null;
  updateVoteButtonState();
  
  try {
    const response = await fetch(`${API_BASE_URL}/votacoes/${votacaoId}/participantes`);
    if (!response.ok) {
      throw new Error('Failed to load participantes');
    }

    const participantes = await response.json();
    renderParticipantes(participantes);
  } catch (error) {
    console.error('Error loading participantes:', error);
    showAlert('Falha ao carregar participantes', 'danger');
    participantsContainer.innerHTML = '<p>Erro ao carregar participantes</p>';
  } finally {
    setLoading(false);
  }
}

// Render participantes
function renderParticipantes(participantes) {
  if (!participantsContainer) return;

  if (participantes.length === 0) {
    participantsContainer.innerHTML = '<p>Nenhum participante nesta votação</p>';
    return;
  }

  let html = '';
  
  participantes.forEach(participante => {
    html += `
      <div class="participant-card" data-id="${participante.id}" onclick="selectParticipante(${participante.id})">
        <img src="${participante.urlFoto}" alt="${participante.nome}" class="participant-image">
        <div class="participant-info">
          <div class="participant-name">${participante.nome}</div>
        </div>
      </div>
    `;
  });
  
  participantsContainer.innerHTML = html;
}

// Select a participante
function selectParticipante(participanteId) {
  // Remove selected class from all participant cards
  document.querySelectorAll('.participant-card').forEach(card => {
    card.classList.remove('selected');
  });
  
  // Add selected class to the clicked participant card
  const selectedCard = document.querySelector(`.participant-card[data-id="${participanteId}"]`);
  if (selectedCard) {
    selectedCard.classList.add('selected');
    selectedParticipanteId = participanteId;
  } else {
    selectedParticipanteId = null;
  }
  
  updateVoteButtonState();
}

// Update vote button state
function updateVoteButtonState() {
  if (!voteButton) return;
  
  voteButton.disabled = !currentVotacaoId || !selectedParticipanteId;
}

// Submit vote
async function submitVote() {
  if (!currentVotacaoId || !selectedParticipanteId) {
    showAlert('Selecione uma votação e um participante', 'danger');
    return;
  }
  
  setLoading(true);
  
  const voteData = {
    participanteId: selectedParticipanteId,
    votacaoId: currentVotacaoId
  };
  
  try {
    const response = await fetch(`${API_BASE_URL}/votos`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(voteData)
    });
    
    if (!response.ok) {
      throw new Error('Failed to submit vote');
    }
    
    // Redirect to success page
    window.location.href = `success.html?votacaoId=${currentVotacaoId}`;
  } catch (error) {
    console.error('Error submitting vote:', error);
    showAlert('Falha ao registrar voto', 'danger');
  } finally {
    setLoading(false);
  }
}

// Make functions available globally
window.selectParticipante = selectParticipante;
