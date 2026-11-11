# Decisions

One short entry per decision: what we chose, what we rejected, why.

## 1. Modular monolith with hexagonal layout
Chose: one app, two modules (identity, board) under internal/, talking only through ports.
Rejected: microservices (too much work for this project), one flat package (no boundaries).
Why: we get clean boundaries without the cost of running many services.
