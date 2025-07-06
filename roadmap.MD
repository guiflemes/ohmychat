## Refactor
- [ ] Improve API ergonomics  
  - Avoid having the user send messages through the output channel directly  
  - Create a method to allow users to send messages properly

- [ ] Create a dedicated context for the chat  
  - Encapsulate all session-related data and logic within this new context

- [ ] Refactor session management  
  - Move session-related components into the new chat context structure

- [ ] Decouple CLI from the bot at initialization  
  - Ensure the CLI does not depend on the bot to start  
  - Pass the context as a parameter instead of having the CLI create or rely on it

## Feat

## BUGS

## TESTS
