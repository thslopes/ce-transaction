# CE Transactions

## Visao Geral

Esta aplicacao tem como objetivo orquestrar a aprovacao de transacoes financeiras com foco em simplicidade e seguranca.

O fluxo de negocio envolve tres grupos principais:

- solicitantes
- aprovadores
- executores

No MVP, a aplicacao nao executa a transacao financeira em si. Ela recebe a solicitacao, conduz o fluxo de aprovacao, registra auditoria, notifica os envolvidos e despacha a solicitacao aprovada para o grupo executor.

## Direcao Arquitetural

A arquitetura recomendada e serverless, orientada a eventos, com implantacao inicial em GCP.

A aplicacao deve ser composta por:

- uma API HTTP pequena para receber solicitacoes e decisoes de aprovacao
- uma camada de persistencia para manter o estado do workflow
- uma camada de eventos para desacoplar notificacoes e processamento assincrono
- um canal de notificacao web com Firebase
- endpoints de relatorio operacional e gerencial

Para o MVP, a entrega para o usuario deve priorizar web responsiva ou PWA. App nativo para iOS fica fora do escopo inicial por aumentar o custo operacional sem trazer beneficio proporcional neste momento.

## Principios do MVP

Os principios obrigatorios do MVP sao:

- simplicidade operacional
- seguranca por padrao
- dupla aprovacao
- trilha completa de auditoria
- idempotencia em operacoes criticas
- relatorios gerenciais basicos

Recomendacao adicional desde o inicio:

- RBAC por perfil e alcada de aprovacao

## Escopo do MVP

Incluido no MVP:

- criacao de solicitacoes financeiras
- fluxo de aprovacao com dois aprovadores
- rejeicao de solicitacoes
- despacho para executores apos aprovacao completa
- notificacoes web com Firebase
- endpoints de relatorio de status e operacao
- auditoria de eventos do processo

Fora do escopo inicial:

- app nativo iOS
- aprovacao por push notification
- execucao financeira pela propria aplicacao
- engine BPM complexa
- automacoes avancadas de SLA e escalonamento

## Plataforma Recomendada

A base tecnica recomendada em GCP e:

- Cloud Run para a API HTTP
- Pub/Sub para eventos de dominio
- Cloud Tasks para retentativas e tarefas assincronas
- Secret Manager para segredos e credenciais
- Cloud Logging para observabilidade e trilhas operacionais
- Firebase Cloud Messaging para notificacoes web

Persistencia recomendada para o inicio:

- Firestore

Armazenamento complementar recomendado:

- Cloud Storage para anexos, exportacoes e consolidados

Alternativa caso o sistema precise cedo de consultas relacionais mais fortes, joins complexos ou transacoes mais rigidas:

- Cloud SQL

## Fluxo de Alto Nivel

1. Um solicitante cria uma solicitacao financeira.
2. A solicitacao entra no estado de pendente de aprovacao.
3. Os aprovadores analisam a solicitacao.
4. Quando a regra de aprovacao e satisfeita, a solicitacao muda para aprovada.
5. A aplicacao despacha a solicitacao para o grupo executor.
6. O sistema envia notificacoes web ao longo do processo.
7. Todos os eventos relevantes ficam registrados para auditoria e relatorios.

## Papeis de Acesso

### Solicitante

Pode:

- criar solicitacoes
- consultar solicitacoes proprias
- acompanhar status

Nao pode:

- aprovar
- rejeitar
- executar

### Aprovador

Pode:

- consultar solicitacoes elegiveis para sua alcada
- aprovar
- rejeitar

Nao pode:

- executar transacoes
- agir fora de sua alcada

### Executor

Pode:

- consultar solicitacoes totalmente aprovadas
- receber solicitacoes despachadas para execucao
- atualizar o andamento operacional, se esse recurso entrar no MVP

Nao pode:

- aprovar ou rejeitar solicitacoes

## Modelo de Workflow

Estados minimos recomendados:

- draft
- pending_approval
- partially_approved
- approved
- rejected
- dispatched
- cancelled

Eventos minimos recomendados:

- RequestCreated
- ApprovalRequested
- ApprovalGranted
- ApprovalRejected
- FullyApproved
- ExecutionDispatched
- NotificationFailed

## Notificacoes no MVP

No MVP, a aplicacao usara notificacoes web com Firebase Cloud Messaging como canal principal de notificacao.

A entrega deve priorizar:

- Android
- desktop web
- PWA instalada

Suporte em iOS web pode existir, mas nao deve ser tratado como base principal do MVP por restricoes de plataforma e menor previsibilidade operacional.

Exemplos de eventos de notificacao:

- nova solicitacao criada
- solicitacao aguardando aprovacao
- solicitacao aprovada
- solicitacao rejeitada
- solicitacao enviada para execucao

As decisoes de aprovar ou rejeitar devem ocorrer em fluxo autenticado da aplicacao, nao diretamente na notificacao.

Fallback recomendado para eventos criticos:

- e-mail

## Endpoints Esperados

A superficie minima da API deve contemplar:

- criacao de solicitacao
- listagem de solicitacoes
- consulta de detalhes da solicitacao
- aprovacao
- rejeicao
- despacho para execucao
- relatorios

Todos os endpoints de escrita devem exigir:

- autenticacao
- autorizacao por papel
- controle de idempotencia

## Relatorios do MVP

Os relatorios iniciais devem cobrir:

- solicitacoes pendentes de aprovacao
- solicitacoes aprovadas
- solicitacoes rejeitadas
- backlog de execucao
- tempo medio de aprovacao
- volume por periodo

## Requisitos de Seguranca

Controles minimos:

- autenticacao forte
- autorizacao por papel
- alcada por valor
- dupla aprovacao
- trilha imutavel de auditoria
- idempotencia
- protecao contra replay
- segregacao de segredos

Identidade recomendada:

- Identity Platform ou integracao com IdP corporativo via OIDC

Os tokens devem carregar, no minimo:

- papel do usuario
- nivel de alcada
- identificador do usuario

## Estrategia de Validacao

A solucao deve ser validada com foco em:

- transicoes corretas de estado
- dupla aprovacao
- prevencao de duplicidade
- autorizacao por papel
- despacho apenas apos aprovacao completa
- retentativas em falhas de notificacao
- rastreabilidade completa para auditoria

## Fases de Implementacao

### Fase 1: Dominio e Seguranca

- definir workflow de negocio
- definir regras de aprovacao
- definir papeis e alcadas
- definir eventos auditaveis

### Fase 2: Base Tecnica

- configurar stack serverless em GCP
- definir persistencia
- definir contratos de evento
- definir notificacoes web com Firebase

### Fase 3: API e Workflow

- implementar endpoints principais
- implementar estados e transicoes
- implementar dupla aprovacao
- implementar idempotencia

### Fase 4: Notificacao e Relatorios

- enviar notificacoes web com Firebase
- expor relatorios operacionais
- consolidar observabilidade

### Fase 5: Hardening

- revisar seguranca
- revisar auditoria
- revisar recuperacao e replay
- preparar operacao

## Proximos Passos

Os proximos refinamentos recomendados sao:

1. detalhar o modelo de dados e o state machine
2. detalhar os endpoints e contratos de evento
3. detalhar a arquitetura GCP com fluxos entre componentes

