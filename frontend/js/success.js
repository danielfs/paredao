// Base URL for API requests
const API_BASE_URL = 'http://localhost:8080';

// DOM Elements
let votacaoTitle;
let totalVotesElement;
let participantsContainer;
let loadingIndicator;
let successContent;

// Current state
let votacaoId = null;
let totalVotes = 0;
let participantesData = [];

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
  // Get DOM elements
  votacaoTitle = document.getElementById('votacao-title');
  totalVotesElement = document.getElementById('total-votes');
  participantsContainer = document.getElementById('participants-container');
  loadingIndicator = document.getElementById('loading-indicator');
  successContent = document.getElementById('success-content');

  // Get votacao ID from URL parameters
  const urlParams = new URLSearchParams(window.location.search);
  votacaoId = urlParams.get('votacaoId');

  if (!votacaoId) {
    showError('Votação não encontrada. Redirecionando para a página inicial...');
    setTimeout(() => {
      window.location.href = 'index.html';
    }, 3000);
    return;
  }

  // Load data
  loadData();
});

// Show error message
function showError(message) {
  const alertContainer = document.getElementById('alert-container');
  alertContainer.innerHTML = `
    <div class="alert alert-danger">
      ${message}
    </div>
  `;
}

// Show/hide loading indicator
function setLoading(isLoading) {
  if (loadingIndicator) {
    loadingIndicator.style.display = isLoading ? 'block' : 'none';
  }
  if (successContent) {
    successContent.style.display = isLoading ? 'none' : 'block';
  }
}

// Load all required data
async function loadData() {
  setLoading(true);
  
  try {
    // Load votacao details
    const votacao = await fetchVotacao(votacaoId);
    
    // Load total votes
    const votacaoTotal = await fetchVotacaoTotal(votacaoId);
    totalVotes = votacaoTotal.total;
    
    // Load votes by participant (now includes all participants, even those with zero votes)
    participantesData = await fetchVotacaoTotalByParticipante(votacaoId);
    
    // Render the data
    renderData(votacao);
  } catch (error) {
    console.error('Error loading data:', error);
    showError('Falha ao carregar dados da votação');
  } finally {
    setLoading(false);
  }
}

// Fetch votacao details
async function fetchVotacao(id) {
  const response = await fetch(`${API_BASE_URL}/votacoes/${id}`);
  if (!response.ok) {
    throw new Error('Failed to load votacao');
  }
  return await response.json();
}

// Fetch total votes for a votacao
async function fetchVotacaoTotal(id) {
  const response = await fetch(`${API_BASE_URL}/estatisticas/votacoes/${id}/total`);
  if (!response.ok) {
    throw new Error('Failed to load total votes');
  }
  return await response.json();
}

// Fetch votes by participant for a votacao
async function fetchVotacaoTotalByParticipante(id) {
  const response = await fetch(`${API_BASE_URL}/estatisticas/votacoes/${id}/participantes`);
  if (!response.ok) {
    throw new Error('Failed to load votes by participant');
  }
  return await response.json();
}

// Fetch all participants
async function fetchAllParticipantes() {
  const response = await fetch(`${API_BASE_URL}/participantes`);
  if (!response.ok) {
    throw new Error('Failed to load all participants');
  }
  return await response.json();
}

// Fetch participants for a specific votacao
async function fetchVotacaoParticipantes(id) {
  const response = await fetch(`${API_BASE_URL}/votacoes/${id}/participantes`);
  if (!response.ok) {
    throw new Error('Failed to load votacao participants');
  }
  return await response.json();
}

// Render all data
function renderData(votacao) {
  // Set votacao title
  if (votacaoTitle) {
    votacaoTitle.textContent = votacao.Descricao;
  }
  
  // Set total votes
  if (totalVotesElement) {
    totalVotesElement.textContent = totalVotes;
  }
  
  // Render participants with vote counts and percentages
  renderParticipants();
}

// Render participants with vote counts and percentages
function renderParticipants() {
  if (!participantsContainer) return;
  
  if (participantesData.length === 0) {
    participantsContainer.innerHTML = '<p>Nenhum participante nesta votação</p>';
    return;
  }
  
  let html = '<div class="participants-grid">';
  
  // Sort participants by vote count (descending)
  participantesData.sort((a, b) => b.total - a.total);
  
  participantesData.forEach((participante, index) => {
    const percentage = totalVotes > 0 ? ((participante.total / totalVotes) * 100).toFixed(1) : 0;
    
    // Assign different colors based on position
    const colorClass = index === 0 ? 'participant-first' : 
                      index === 1 ? 'participant-second' : 
                      index === 2 ? 'participant-third' : 'participant-other';
    
    // Use a placeholder image since the API doesn't provide url_foto in this endpoint
    const placeholderImage = `https://ui-avatars.com/api/?name=${encodeURIComponent(participante.nome)}&background=random&color=fff&size=100`;
    
    html += `
      <div class="participant-vertical ${colorClass}">
        <h3 class="participant-name">${participante.nome}</h3>
        <div class="participant-image-container">
          <img src="${placeholderImage}" alt="${participante.nome}" class="participant-photo">
        </div>
        <div class="participant-percentage-container">
          <span class="participant-percentage">${percentage}%</span>
          <span class="votes-count">(${participante.total} votos)</span>
        </div>
      </div>
    `;
  });
  
  html += '</div>';
  participantsContainer.innerHTML = html;
}
