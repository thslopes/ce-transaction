# CE Transactions

## Visao Geral

Esta aplicacao tem como objetivo orquestrar a aprovacao de transacoes financeiras com foco em simplicidade e seguranca.

O fluxo de negocio envolve tres grupos principais:

- solicitantes
- aprovadores
- executores

O sistema tambem possui um perfil administrativo para gerenciar o ciclo de acesso dos usuarios.

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

No MVP, os perfis sao uma propriedade do proprio usuario.

## Escopo do MVP

Incluido no MVP:

- cadastro interno de usuarios por link temporario
- validacao administrativa de usuarios
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

Uso sugerido da persistencia operacional:

- usuarios com status `active` ou `inactive`
- perfis associados ao usuario
- solicitacoes financeiras e suas transicoes
- eventos de auditoria

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

### Admin

Pode:

- gerar link de cadastro
- listar usuarios
- validar usuario associando perfis e ativando acesso
- desativar usuario
- consultar trilha administrativa

Nao pode automaticamente:

- substituir regras de negocio de aprovacao sem possuir tambem o perfil adequado

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

- geracao de link de cadastro
- cadastro por token valido
- listagem de usuarios
- validacao de usuario
- desativacao de usuario
- consulta de perfil do usuario autenticado
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

Observacoes sobre acesso:

- concluir cadastro nao concede acesso automaticamente
- acesso ao sistema exige usuario `active`
- acesso ao sistema exige ao menos um perfil associado

## Telas Esperadas no MVP

A aplicacao deve expor ao menos as seguintes telas ou paginas:

- login e acesso autenticado
- listagem de usuarios para admin
- validacao de usuario para admin
- dashboard inicial por perfil
- criacao de solicitacao financeira
- listagem de solicitacoes
- detalhe de solicitacao
- fila de aprovacoes para aprovadores
- fila de execucao para executores
- relatorios operacionais
- preferencias de notificacao

Resumo funcional por tela:

- Dashboard: exibir pendencias, itens recentes e atalhos conforme o perfil do usuario.
- Gestao de usuarios: permitir listar usuarios, inspecionar status, perfis e ativar acesso.
- Validacao de usuario: permitir ao admin definir perfis e ativar o usuario em uma unica acao.
- Criacao de solicitacao: permitir abertura de nova solicitacao com dados obrigatorios e anexos, quando houver.
- Listagem: permitir filtro por status, periodo, solicitante e responsavel.
- Detalhe: mostrar historico, aprovacoes, auditoria e dados completos da solicitacao.
- Fila de aprovacoes: destacar itens pendentes de decisao dentro da alcada do aprovador.
- Fila de execucao: mostrar apenas itens totalmente aprovados e prontos para despacho ou acompanhamento operacional.
- Relatorios: consolidar volume, tempo medio, pendencias e backlog.
- Preferencias de notificacao: permitir controle de permissao, dispositivo e fallback de notificacao.

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

- cadastro interno no MVP, com evolucao futura para Identity Platform ou integracao com IdP corporativo via OIDC

No MVP, o controle de acesso deve considerar, no minimo:

- identificador do usuario
- status `active` ou `inactive`
- perfis do usuario
- nivel de alcada, quando aplicavel para aprovacao

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

- definir fluxo de cadastro e validacao administrativa
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

## Documentacao Complementar

- Casos de uso do MVP: [docs/use-cases.md](/home/thiago-slopes/workspace/ce-transactions/docs/use-cases.md)

