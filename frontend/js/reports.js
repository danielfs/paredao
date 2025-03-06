// Base URL for API requests
const REPORTS_API_BASE_URL = 'http://localhost:8080';

// DOM Elements
let reportVotacaoSelect;
let reportContent;
let totalVotesElement;
let votesByParticipantElement;
let votesByHourElement;

// Initialize the reports functionality
document.addEventListener('DOMContentLoaded', () => {
  console.log('Reports: DOM Content Loaded');
  
  // Get DOM elements
  reportVotacaoSelect = document.getElementById('report-votacao-select');
  reportContent = document.getElementById('report-content');
  totalVotesElement = document.getElementById('total-votes');
  votesByParticipantElement = document.getElementById('votes-by-participant');
  votesByHourElement = document.getElementById('votes-by-hour');
  
  console.log('Reports: DOM Elements', {
    reportVotacaoSelect,
    reportContent,
    totalVotesElement,
    votesByParticipantElement,
    votesByHourElement
  });

  // Set up event listeners
  setupReportsEventListeners();

  // Load votacoes for the dropdown
  loadVotacoesForReports();
});

// Set up event listeners for reports
function setupReportsEventListeners() {
  console.log('Reports: Setting up event listeners');
  
  // When a votacao is selected, load its reports
  reportVotacaoSelect.addEventListener('change', () => {
    console.log('Reports: Votacao selected', reportVotacaoSelect.value);
    
    const votacaoId = reportVotacaoSelect.value;
    if (votacaoId) {
      loadReportsForVotacao(votacaoId);
    } else {
      reportContent.style.display = 'none';
    }
  });
}

// Load votacoes for the reports dropdown
async function loadVotacoesForReports() {
  console.log('Reports: Loading votacoes for dropdown');
  
  try {
    const response = await fetch(`${REPORTS_API_BASE_URL}/votacoes`);
    if (!response.ok) {
      throw new Error('Failed to load votacoes');
    }

    const votacoes = await response.json();
    console.log('Reports: Votacoes loaded', votacoes);
    
    populateVotacoesDropdown(votacoes);
  } catch (error) {
    console.error('Error loading votacoes for reports:', error);
    // Use console.error instead of showAlert since showAlert is defined in admin.js
    console.error('Failed to load votacoes for reports');
  }
}

// Populate the votacoes dropdown
function populateVotacoesDropdown(votacoes) {
  console.log('Reports: Populating dropdown with votacoes', votacoes);
  
  reportVotacaoSelect.innerHTML = '<option value="">Selecione uma votação</option>';
  
    votacoes.forEach(votacao => {
    reportVotacaoSelect.innerHTML += `<option value="${votacao.id}">${votacao.descricao}</option>`;
  });
  
  console.log('Reports: Dropdown populated');
}

// Load reports for a selected votacao
async function loadReportsForVotacao(votacaoId) {
  console.log('Reports: Loading reports for votacao', votacaoId);
  
  try {
    // Show loading state
    totalVotesElement.innerHTML = '<div class="spinner"></div><p>Carregando...</p>';
    votesByParticipantElement.innerHTML = '<div class="spinner"></div><p>Carregando...</p>';
    votesByHourElement.innerHTML = '<div class="spinner"></div><p>Carregando...</p>';
    
    // Show the report content
    reportContent.style.display = 'block';
    
    console.log('Reports: Fetching data from API');
    
    // Load all three reports in parallel
    const [totalResponse, participantResponse, hourlyResponse] = await Promise.all([
      fetch(`${REPORTS_API_BASE_URL}/estatisticas/votacoes/${votacaoId}/total`),
      fetch(`${REPORTS_API_BASE_URL}/estatisticas/votacoes/${votacaoId}/participantes`),
      fetch(`${REPORTS_API_BASE_URL}/estatisticas/votacoes/${votacaoId}/hourly`)
    ]);
    
    // Check responses
    if (!totalResponse.ok || !participantResponse.ok || !hourlyResponse.ok) {
      throw new Error('Failed to load one or more reports');
    }
    
    // Parse JSON responses
    const totalData = await totalResponse.json();
    const participantData = await participantResponse.json();
    const hourlyData = await hourlyResponse.json();
    
    console.log('Reports: Data received', {
      totalData,
      participantData,
      hourlyData
    });
    
    // Display the reports
    displayTotalVotes(totalData);
    displayVotesByParticipant(participantData);
    displayVotesByHour(hourlyData);
  } catch (error) {
    console.error('Error loading reports:', error);
    reportContent.style.display = 'block';
    totalVotesElement.innerHTML = '<p class="error">Erro ao carregar dados</p>';
    votesByParticipantElement.innerHTML = '<p class="error">Erro ao carregar dados</p>';
    votesByHourElement.innerHTML = '<p class="error">Erro ao carregar dados</p>';
  }
}

// Display total votes
function displayTotalVotes(data) {
  console.log('Reports: Displaying total votes', data);
  
  totalVotesElement.innerHTML = `
    <div class="total-count">
      <span class="count-number">${data.total}</span>
      <span class="count-label">votos totais</span>
    </div>
  `;
}

// Display votes by participant
function displayVotesByParticipant(data) {
  console.log('Reports: Displaying votes by participant', data);
  
  if (data.length === 0) {
    votesByParticipantElement.innerHTML = '<p>Nenhum dado disponível</p>';
    return;
  }
  
  let html = `
    <table class="report-table">
      <thead>
        <tr>
          <th>Participante</th>
          <th>Total de Votos</th>
          <th>Porcentagem</th>
        </tr>
      </thead>
      <tbody>
  `;
  
  // Calculate total votes for percentage calculation
  const totalVotes = data.reduce((sum, item) => sum + item.total, 0);
  
  data.forEach(item => {
    const percentage = totalVotes > 0 ? ((item.total / totalVotes) * 100).toFixed(2) : '0.00';
    
    html += `
      <tr>
        <td>${item.nome}</td>
        <td>${item.total}</td>
        <td>${percentage}%</td>
      </tr>
    `;
  });
  
  html += `
      </tbody>
    </table>
  `;
  
  votesByParticipantElement.innerHTML = html;
}

// Display votes by hour
function displayVotesByHour(data) {
  console.log('Reports: Displaying votes by hour', data);
  
  if (data.length === 0) {
    votesByHourElement.innerHTML = '<p>Nenhum dado disponível</p>';
    return;
  }
  
  let html = `
    <table class="report-table">
      <thead>
        <tr>
          <th>Hora</th>
          <th>Total de Votos</th>
        </tr>
      </thead>
      <tbody>
  `;
  
  data.forEach(item => {
    const hourLabel = formatHour(item.hour);
    
    html += `
      <tr>
        <td>${hourLabel}</td>
        <td>${item.total}</td>
      </tr>
    `;
  });
  
  html += `
      </tbody>
    </table>
  `;
  
  votesByHourElement.innerHTML = html;
}

// Format hour for display (24-hour format)
function formatHour(hour) {
  return `${hour.toString().padStart(2, '0')}:00`;
}
