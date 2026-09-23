# Casos de Uso do MVP

## Objetivo

Este documento detalha os principais casos de uso do MVP da aplicacao de aprovacao de transacoes financeiras. O foco e transformar a visao arquitetural em comportamento verificavel, servindo como base para testes, modelagem de dominio e implementacao incremental.

## Atores

- Admin: gerencia usuarios, perfis e status de acesso.
- Candidato: utiliza o link de cadastro para criar sua conta inicial.
- Solicitante: cria e acompanha solicitacoes.
- Aprovador: analisa e decide sobre solicitacoes dentro de sua alcada.
- Executor: recebe solicitacoes totalmente aprovadas para processamento operacional.
- Sistema: aplica regras de negocio, registra auditoria e dispara notificacoes.

## Convencoes

- Toda acao relevante deve gerar evento de auditoria.
- Toda operacao de escrita deve ser idempotente quando aplicavel.
- Aprovacoes fora da alcada devem ser rejeitadas.
- Solicitacoes rejeitadas nao podem voltar para aprovadas sem novo fluxo definido.
- O cadastro do usuario e interno no MVP.
- O cadastro nao concede acesso automatico ao sistema.
- O acesso efetivo exige usuario com `status = active` e ao menos um perfil associado.
- Os perfis sao uma propriedade do usuario.

## UC-00 Gerar Link de Cadastro

### Objetivo

Permitir que o admin gere um link temporario para que um candidato realize seu cadastro inicial.

### Ator principal

- Admin

### Pre-condicoes

- O admin esta autenticado.
- O admin possui perfil `admin`.

### Fluxo principal

1. O admin acessa a area de administracao de usuarios.
2. O admin solicita a geracao de um link de cadastro.
3. O sistema gera um token assinado com validade curta.
4. O sistema monta o link de cadastro com o token.
5. O admin compartilha o link com o candidato.

### Fluxos alternativos

- Se o usuario nao possuir perfil `admin`, a acao e negada.

### Resultado esperado

- Link temporario de cadastro gerado com sucesso.

### Testes minimos derivados

- Admin consegue gerar link temporario.
- Usuario sem perfil `admin` nao consegue gerar link.

## UC-01 Realizar Cadastro

### Objetivo

Permitir que um candidato utilize um link valido para criar sua conta inicial no sistema.

### Ator principal

- Candidato

### Pre-condicoes

- O candidato possui link com token valido.

### Fluxo principal

1. O candidato acessa o link de cadastro.
2. O sistema valida assinatura e expiracao do token.
3. O candidato informa nome, e-mail, telefone e senha.
4. O sistema valida os dados obrigatorios.
5. O sistema cria o usuario com `status = inactive`.
6. O sistema grava o usuario com `profiles = []`.
7. O sistema registra evento de auditoria de criacao do usuario.

### Fluxos alternativos

- Se o token estiver expirado ou invalido, o cadastro e bloqueado.
- Se ja existir usuario com o identificador unico adotado pelo sistema, o cadastro e rejeitado.

### Resultado esperado

- Usuario criado sem acesso efetivo ao sistema.

### Testes minimos derivados

- Cadastro com token valido cria usuario `inactive`.
- Cadastro com token expirado falha.
- Cadastro duplicado para o mesmo identificador falha.

## UC-02 Listar Usuarios

### Objetivo

Permitir que o admin visualize usuarios cadastrados, seus perfis e seu status de acesso.

### Ator principal

- Admin

### Pre-condicoes

- O admin esta autenticado.
- O admin possui perfil `admin`.

### Fluxo principal

1. O admin acessa a listagem de usuarios.
2. O sistema retorna usuarios com nome, e-mail, telefone, status, perfis e datas relevantes.
3. O admin filtra ou ordena os resultados.

### Resultado esperado

- O admin visualiza os usuarios necessarios para validacao e gestao.

### Testes minimos derivados

- Admin consegue listar usuarios.
- Usuario sem perfil `admin` nao consegue listar usuarios.

## UC-03 Validar Usuario

### Objetivo

Permitir que o admin libere o acesso de um usuario cadastrando seus perfis e ativando sua conta em uma unica acao.

### Ator principal

- Admin

### Pre-condicoes

- O admin esta autenticado.
- O admin possui perfil `admin`.
- O usuario alvo existe no sistema.
- O usuario alvo esta `inactive`.

### Fluxo principal

1. O admin acessa o detalhe de um usuario.
2. O admin define um ou mais perfis para o usuario.
3. O admin confirma a validacao.
4. O sistema associa os perfis ao usuario.
5. O sistema altera o `status` para `active`.
6. O sistema registra `validated_at` e `validated_by`.
7. O sistema grava auditoria da validacao.

### Fluxos alternativos

- Se nenhum perfil for informado, a validacao e rejeitada.
- Se o usuario ja estiver ativo, a operacao e tratada como atualizacao ou rejeitada, conforme a regra adotada.

### Resultado esperado

- Usuario ativo e apto a acessar o sistema conforme seus perfis.

### Testes minimos derivados

- Validacao com perfis ativa o usuario.
- Validacao sem perfis falha.
- Usuario ativo com perfis corretos pode acessar rotas protegidas.

## UC-04 Desativar Usuario

### Objetivo

Permitir que o admin suspenda o acesso de um usuario sem apagar seu historico.

### Ator principal

- Admin

### Pre-condicoes

- O admin esta autenticado.
- O admin possui perfil `admin`.
- O usuario alvo existe no sistema.

### Fluxo principal

1. O admin seleciona um usuario.
2. O admin solicita a desativacao.
3. O sistema altera o `status` para `inactive`.
4. O sistema registra auditoria.

### Resultado esperado

- Usuario perde acesso sem remocao do historico.

### Testes minimos derivados

- Usuario desativado nao acessa rotas protegidas.

## UC-05 Autenticar Usuario

### Objetivo

Permitir que um usuario cadastrado acesse a aplicacao com suas credenciais internas e receba o contexto de autorizacao adequado.

### Ator principal

- Usuario

### Pre-condicoes

- O usuario concluiu cadastro.
- O usuario possui credenciais validas.

### Fluxo principal

1. O usuario acessa a tela de login.
2. O usuario informa suas credenciais.
3. O sistema valida as credenciais internas.
4. O sistema carrega o usuario interno.
5. O sistema verifica se o usuario esta `active`.
6. O sistema carrega os perfis do usuario.
7. O sistema cria a sessao autenticada ou emite o token de acesso.
8. O usuario e encaminhado ao dashboard adequado.

### Fluxos alternativos

- Se as credenciais forem invalidas, o acesso e negado.
- Se o usuario estiver `inactive`, o acesso e negado.
- Se o usuario estiver ativo mas sem perfis, o acesso e negado.

### Resultado esperado

- Usuario autenticado com contexto de autorizacao carregado.

### Testes minimos derivados

- Usuario ativo com perfil recebe sessao valida.
- Usuario `inactive` nao acessa o sistema.
- Usuario sem perfil nao acessa o sistema.

## UC-06 Criar Solicitacao Financeira

### Objetivo

Permitir que um solicitante abra uma nova solicitacao financeira para posterior aprovacao.

### Ator principal

- Solicitante

### Pre-condicoes

- O usuario esta autenticado.
- O usuario possui perfil de solicitante.
- Os dados obrigatorios da solicitacao estao disponiveis.

### Fluxo principal

1. O solicitante acessa a tela de criacao.
2. O solicitante informa os dados obrigatorios da solicitacao.
3. O sistema valida os campos e regras minimas de preenchimento.
4. O sistema cria a solicitacao com status inicial `pending_approval`.
5. O sistema registra evento de auditoria de criacao.
6. O sistema identifica os aprovadores elegiveis.
7. O sistema prepara o envio de notificacoes aos aprovadores.

### Fluxos alternativos

- Se os dados obrigatorios estiverem incompletos, o sistema rejeita a criacao e informa os erros.
- Se o solicitante reenviar a mesma requisicao com a mesma chave de idempotencia, o sistema retorna o mesmo resultado sem duplicar a solicitacao.

### Resultado esperado

- Solicitacao criada com identificador unico.
- Status inicial definido.
- Auditoria registrada.

### Testes minimos derivados

- Criar solicitacao valida resulta em `pending_approval`.
- Criacao duplicada com a mesma chave nao gera nova solicitacao.
- Criacao invalida falha sem persistir estado parcial.

## UC-07 Aprovar Solicitacao Dentro da Alcada

### Objetivo

Permitir que um aprovador registre aprovacao valida de uma solicitacao pendente.

### Ator principal

- Aprovador

### Pre-condicoes

- O usuario esta autenticado.
- O usuario possui perfil de aprovador.
- A solicitacao esta em estado aprovavel.
- O aprovador possui alcada suficiente.

### Fluxo principal

1. O aprovador acessa sua fila de aprovacoes.
2. O aprovador abre o detalhe da solicitacao.
3. O aprovador registra a aprovacao.
4. O sistema valida alcada, estado atual e duplicidade.
5. O sistema registra a aprovacao.
6. O sistema atualiza o status para `partially_approved` ou `approved`, conforme a quantidade de aprovacoes validas.
7. O sistema registra evento de auditoria.

### Fluxos alternativos

- Se o aprovador nao possuir alcada, o sistema rejeita a acao.
- Se a mesma aprovacao for reenviada, o sistema nao conta a aprovacao duas vezes.
- Se a solicitacao ja estiver rejeitada ou finalizada, o sistema rejeita a acao.

### Resultado esperado

- Aprovacao registrada uma unica vez.
- Status recalculado corretamente.

### Testes minimos derivados

- Primeira aprovacao valida move para `partially_approved`.
- Segunda aprovacao valida move para `approved`.
- Aprovacao duplicada nao altera estado.
- Aprovacao fora da alcada falha.

## UC-08 Rejeitar Solicitacao

### Objetivo

Permitir que um aprovador rejeite uma solicitacao dentro de sua responsabilidade.

### Ator principal

- Aprovador

### Pre-condicoes

- O usuario esta autenticado.
- O usuario possui perfil de aprovador.
- A solicitacao esta em estado rejeitavel.

### Fluxo principal

1. O aprovador acessa o detalhe da solicitacao.
2. O aprovador informa a rejeicao.
3. O sistema valida permissao e estado.
4. O sistema altera o status para `rejected`.
5. O sistema registra evento de auditoria.
6. O sistema prepara notificacao para o solicitante.

### Fluxos alternativos

- Se a solicitacao ja estiver finalizada, a rejeicao nao e permitida.
- Se o usuario nao possuir permissao, a acao e negada.

### Resultado esperado

- Solicitacao rejeitada e bloqueada para novas aprovacoes no fluxo atual.

### Testes minimos derivados

- Rejeicao valida move para `rejected`.
- Solicitacao rejeitada nao aceita nova aprovacao.

## UC-09 Consolidar Dupla Aprovacao

### Objetivo

Garantir que a solicitacao so seja considerada aprovada apos cumprir a regra de dupla aprovacao.

### Ator principal

- Sistema

### Pre-condicoes

- Existe uma solicitacao com aprovacoes registradas.
- A regra de dupla aprovacao esta ativa para o fluxo.

### Fluxo principal

1. O sistema avalia a quantidade de aprovacoes validas.
2. O sistema verifica se as aprovacoes atendem as regras configuradas.
3. O sistema mantem `partially_approved` enquanto faltar aprovacao valida.
4. O sistema muda para `approved` quando a regra for satisfeita.
5. O sistema registra evento de auditoria de consolidacao.

### Fluxos alternativos

- Se houver rejeicao previa, a consolidacao nao ocorre.
- Se uma aprovacao for invalida, ela nao entra no calculo.

### Resultado esperado

- Nenhuma solicitacao e aprovada antes da segunda aprovacao valida.

### Testes minimos derivados

- Uma unica aprovacao nao aprova a solicitacao.
- Duas aprovacoes validas aprovam a solicitacao.

## UC-10 Enviar Solicitacao Aprovada para Execucao

### Objetivo

Permitir que uma solicitacao totalmente aprovada seja encaminhada ao grupo executor.

### Ator principal

- Sistema
- Executor

### Pre-condicoes

- A solicitacao esta com status `approved`.

### Fluxo principal

1. O sistema identifica que a solicitacao esta pronta para execucao.
2. O sistema a disponibiliza na fila de execucao.
3. O sistema registra evento de despacho.
4. O sistema prepara notificacao para executores elegiveis.

### Fluxos alternativos

- Se a solicitacao nao estiver `approved`, ela nao vai para a fila de execucao.

### Resultado esperado

- Apenas solicitacoes totalmente aprovadas aparecem para execucao.

### Testes minimos derivados

- Solicitacao aprovada entra na fila de execucao.
- Solicitacao nao aprovada nao entra na fila de execucao.

## UC-11 Consultar Fila de Aprovacao

### Objetivo

Permitir que aprovadores visualizem apenas itens relevantes para sua atuacao.

### Ator principal

- Aprovador

### Pre-condicoes

- O usuario esta autenticado como aprovador.

### Fluxo principal

1. O aprovador acessa a fila de aprovacoes.
2. O sistema lista solicitacoes pendentes compativeis com seu perfil e alcada.
3. O aprovador filtra ou ordena os itens.

### Resultado esperado

- O aprovador visualiza apenas itens elegiveis para decisao.

### Testes minimos derivados

- Itens fora da alcada nao aparecem na fila.
- Itens rejeitados ou concluidos nao aparecem como pendentes.

## UC-12 Consultar Fila de Execucao

### Objetivo

Permitir que executores visualizem somente o que esta pronto para processamento operacional.

### Ator principal

- Executor

### Pre-condicoes

- O usuario esta autenticado como executor.

### Fluxo principal

1. O executor acessa a fila de execucao.
2. O sistema lista apenas solicitacoes com status apto para execucao.
3. O executor consulta o detalhe do item.

### Resultado esperado

- O executor nao visualiza itens ainda pendentes de aprovacao.

### Testes minimos derivados

- Apenas itens aprovados sao retornados para o executor.

## UC-13 Consultar Relatorio Operacional

### Objetivo

Permitir visibilidade gerencial sobre volume, status e backlog do processo.

### Ator principal

- Usuario autorizado

### Pre-condicoes

- O usuario esta autenticado.
- O usuario possui permissao para relatorios.

### Fluxo principal

1. O usuario acessa a area de relatorios.
2. O sistema consolida dados por status, periodo e fila.
3. O sistema exibe totais e indicadores principais.

### Resultado esperado

- Relatorio operacional com pendencias, aprovadas, rejeitadas e backlog.

### Testes minimos derivados

- O relatorio retorna contagens coerentes com os estados das solicitacoes.

## Regras Transversais de Acesso

- Toda rota protegida exige usuario autenticado.
- Toda rota protegida exige usuario interno existente.
- Toda rota protegida exige `status = active`.
- Toda rota protegida exige ao menos um perfil compativel com a operacao.
- Rotas de aprovacao tambem exigem validacao de alcada.

## Ordem Sugerida de Implementacao em Baby-Steps

1. UC-00 Gerar Link de Cadastro
2. UC-01 Realizar Cadastro
3. UC-02 Listar Usuarios
4. UC-03 Validar Usuario
5. UC-05 Autenticar Usuario
6. UC-06 Criar Solicitacao Financeira
7. UC-07 Aprovar Solicitacao Dentro da Alcada
8. UC-09 Consolidar Dupla Aprovacao
9. UC-08 Rejeitar Solicitacao
10. UC-10 Enviar Solicitacao Aprovada para Execucao
11. UC-11 Consultar Fila de Aprovacao
12. UC-12 Consultar Fila de Execucao
13. UC-13 Consultar Relatorio Operacional

## Primeiro Slice Recomendado

O primeiro slice recomendado para implementacao e testes de acesso e administracao e:

- admin gera link temporario de cadastro
- candidato conclui cadastro com nome, e-mail, telefone e senha
- usuario nasce `inactive`
- admin valida usuario em acao unica, definindo perfis e ativando acesso
- usuario ativo com perfil consegue autenticar
- usuario sem perfil ou `inactive` nao acessa o sistema

O primeiro slice recomendado para implementacao e testes de negocio, depois disso, e:

- usuario autenticado com papel carregado
- criar solicitacao valida
- registrar primeira aprovacao valida
- registrar segunda aprovacao valida
- consolidar estado final como `approved`
- impedir aprovacao duplicada
- impedir aprovacao fora da alcada

## Documentacao Complementar

- Contratos de acesso e administracao: [docs/access-api.md](/home/thiago-slopes/workspace/ce-transactions/docs/access-api.md)