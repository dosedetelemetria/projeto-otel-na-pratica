# Guia de Contribuição

## **Resumo**

Este guia orienta como contribuir ao projeto de forma organizada e colaborativa. Siga as etapas para validar ou criar Issues, alinhar responsabilidades e abrir Pull Requests, garantindo que o trabalho seja organizado e eficiente.

## **Sumário**

- [Fluxo de Contribuição](#fluxo-de-contribuição)
- [Abrir um Pull Request](#abrir-um-pull-request)

---

## **Fluxo de Contribuição**

> ⚠️ **Importante:** Toda contribuição deve começar com a criação ou validação de uma Issue.

1. **Verifique as Issues Existentes**:
   - Antes de começar, cheque a aba **Issues** para evitar duplicações.
   - Caso encontre uma Issue relacionada, comente nela para alinhar com o responsável.

2. **Crie uma Nova Issue**:
   - Se não encontrar algo relacionado, clique em **New Issue** e preencha:
     - **Título:** Um resumo claro do problema ou melhoria.
     - **Descrição:** Contexto, solução proposta e impacto esperado.

3. **Atualize a Issue**:
   - Comente na Issue para informar que está trabalhando nela.
   - Adicione-se como responsável (em **Assignees**) para sinalizar que está cuidando do item.

4. **Encerrar a Issue**:
   - Ao finalizar, abra um Pull Request (PR) que referencie a Issue:
     - Exemplo: _"Este PR resolve a Issue #1."_  
   - A Issue será fechada automaticamente quando o PR for mesclado.

---

## **Abrir um Pull Request**

Siga os passos para enviar suas contribuições:

1. **Fork do repositório:**  
   No GitHub, clique em **"Fork"** para criar uma cópia do repositório.

2. **Clone do fork para sua máquina local:**
   ```bash
   git clone https://github.com/seu-usuario/projeto-otel-na-pratica.git
   cd projeto-otel-na-pratica
   ```

3. **Sincronize com o repositório principal:**
   ```bash
   git remote add upstream https://github.com/dosedetelemetria/projeto-otel-na-pratica.git
   ```

4. **Crie uma branch para suas mudanças:**
   ```bash
   git checkout -b minha-branch
   ```

5. **Após as mudanças, envie para o fork:**
   ```bash
   git add .
   git commit -m "Descrição do que foi feito" -m "Se necessario, detalhe o que foi feito e o que o seu commit resolve"
   git push origin minha-branch
   ```

6. **Crie o Pull Request no GitHub:**

   - Vá para o repositório do seu fork no GitHub.
   - Você verá um botão "Compare & Pull Request" após enviar a branch.
   - Clique no botão e preencha a descrição conforme solicitado.
   - Clique em "Create Pull Request" para enviar.
   - Aguarde a revisão do Pull Request.
   - Após aprovado você está apto para mesclar o Pull Request.
