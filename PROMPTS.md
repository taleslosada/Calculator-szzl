# Prompts used

I used **Claude** (Anthropic) in a single conversation, in Spanish, to translate a Java calculator into Python and Go.
Below are my prompts, translated to English, in the order I sent them.

---

### 1. Kick-off

> I'm applying for a job interview and they left me a take-home assignment lets say: *(pasted the full assignment: objective, functional and non-functional requirements, constraints, deliverables, and instructions)*.
> Tell me what you think. I already have a calculator in Java from a couple of years back that i nearly deleted. Is it worth reusing patterns from this previous project and translated?

**Outcome:** Claude suggested reusing the *approach* from that project (domain separated from transport, table-driven tests, coverage) and requested to see the code from the previous Java calculator. The proposed stack whit all the changes made to the previous files was a Go standard-library backend, React + TypeScript + Vite, Vitest + Testing Library, a multi-stage Dockerfile, and CI.

### 2. Scoping and review plan

> Before anything, tell me: 1) how complex you think it is, 2) which files I should review most carefully, and 3) anything else you briefly recommend.

**Outcome:** The assessment was low-to-medium complexity, with the real bar being code quality, edge cases and tests. Files to review closely:

- `calculator.go` (sentinel errors, NaN/Inf guards)
- `handler.go` (400 vs 422 mapping)
- the table-driven tests
- `client.ts`
- `useCalculator.ts`

The recommendations were small meaningful commits, transparency about AI usage and fo course running everything locally.

### 3. Build-up

> Okay, lets get to it.

**Outcome:** Claude proceeded to translate the backend, update the frontend, make tests, and build a Dockerfile. It then ran `go vet`, `go test -cover`, `npm run lint`, `npm run coverage` and `npm run build`, smoke-tested the API with `curl`, and checked the UI in a headless browser at desktop and mobile widths.

---

## How I worked with the output

- I reviewed each layer, focusing on the files listed above, and made sure I can explain every design decision in the README.
- I ran the full test suite, the app and the Docker image on my own machine before submitting.
