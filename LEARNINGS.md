
## 2025-10-26 — PR #2 Howard-nolan/go-http-server
- **Why:** Keeping a running record of lessons learned makes progress visible and reinforces key concepts from each phase of the project.  
- **How:** Used a GitHub Action triggered on merged pull requests to extract the “What I Learned” section and append it to .  
- **Gotcha:** Initially forgot that Actions only trigger from the default branch; ensured the workflow YAML lives on .  
- **Next time:** Test new workflows with a small dummy PR like this one before relying on them in production.

