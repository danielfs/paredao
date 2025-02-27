-- Create database if it doesn't exist
CREATE DATABASE IF NOT EXISTS paredao;

-- Use the database
USE paredao;

-- Create participantes table
CREATE TABLE IF NOT EXISTS participantes (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    url_foto VARCHAR(255) NOT NULL
);

-- Create votacoes table
CREATE TABLE IF NOT EXISTS votacoes (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    descricao VARCHAR(255) NOT NULL
);

-- Create votacao_participante table
CREATE TABLE IF NOT EXISTS votacao_participante (
    participante_id BIGINT NOT NULL,
    votacao_id BIGINT NOT NULL,
    PRIMARY KEY (participante_id, votacao_id),
    FOREIGN KEY (participante_id) REFERENCES participantes(id) ON DELETE CASCADE,
    FOREIGN KEY (votacao_id) REFERENCES votacoes(id) ON DELETE CASCADE
);


-- Create votos table
CREATE TABLE IF NOT EXISTS votos (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    participante_id BIGINT NOT NULL,
    votacao_id BIGINT NOT NULL,
    data_hora DATETIME NOT NULL,
    FOREIGN KEY (participante_id) REFERENCES participantes(id) ON DELETE CASCADE,
    FOREIGN KEY (votacao_id) REFERENCES votacoes(id) ON DELETE CASCADE
);


-- Insert data
SET NAMES utf8mb4;

INSERT INTO `participantes` (`id`, `nome`, `url_foto`) VALUES
(1,	'Johann Sebastian Bach',	'https://hips.hearstapps.com/hmg-prod/images/johann-sebastian-bach-gettyimages-51246888.jpg'),
(2,	'Antonio Lucio Vivaldi',	'https://classicosdosclassicos.mus.br/wp-content/uploads/2021/03/antonio_vivaldi.png'),
(3,	'Beethoven',	'https://i.natgeofe.com/n/42e08d5a-5fbd-4c02-aa50-7e8a237aea72/16-beethoven-portrait-og_square.jpg');

INSERT INTO `votacoes` (`id`, `descricao`) VALUES
(1,	'Votação 27/02/2025');

INSERT INTO `votacao_participante` (`participante_id`, `votacao_id`) VALUES
(1,	1),
(2,	1),
(3,	1);
