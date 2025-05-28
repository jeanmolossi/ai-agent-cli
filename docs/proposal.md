# Proposta

## Do produto

O produto é um agente de IA que auxilia no desenvolvimento de software, onde o dev
pode fazer perguntas/pedidos de tarefas ou recomendações para o agente e ele terá
o contexto todo do projeto que está sendo desenvolvido para conseguir auxiliar de
forma assertiva.

É importante que esse agente possa integrar com diferentes LLM com um "contrato"
padrão e também possa rodar local, utilizando de LLMs locais para tal processamento.

## Das tecnologias

Vamos desenvolver utilizando Golang para otimizar ao máximo o processamento.
Também vamos utilizar sqlite como banco de dados auxiliar quando necessário para
manter informações importantes para o agente funcionar.

Inicialmente teremos uma CLI (command line interface), porém a estrutura deve ser
projetada imaginando que isso possa vir a ter uma API, um Action do Github e outros
tipos de integrações externas.

## Da arquitetura

Esse agente deverá ser "executável" em qualquer ambiente:

- Local;
- Nuvem;
- Serverless;
- Containers;

A organização do projeto também é importante, é preferível utilizar designs escaláveis
e de fácil manutenção.
