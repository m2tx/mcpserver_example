{{.Header}}

The search results below were collected via multiple targeted queries in English and Portuguese. Analyze them carefully and produce a structured adverse media report.

---
{{.SearchResults}}
---

## Instructions

### Step 1 — Name Disambiguation

Before classifying any finding, determine whether it refers to **{{.Name}}** or a different person who shares the same name. Apply the following rules:

- **Confirmed**: the full name (first and last) appears together AND at least one additional identifier corroborates the match (organization, country, role, date, or photo caption).
- **Probable**: the full name matches exactly but no additional identifier is present.
- **Excluded**: only a partial name matches (same last name only, same first name only, or initials) with no corroborating identifiers.

### Step 2 — Classify Each Finding

Assign one or more of the following categories:

| Code | Category | Examples |
|------|----------|---------|
| FC | Financial Crime | fraud, money laundering, embezzlement, tax evasion, Ponzi scheme |
| CB | Corruption / Bribery | bribery, kickbacks, public corruption |
| LC | Legal / Criminal | arrest, indictment, conviction, lawsuit, homicide, kidnapping, human trafficking, sexual assault, terrorism, torture, genocide, robbery |
| RS | Regulatory / Sanctions | sanctions lists, regulatory actions, license revocations, debarment |
| RP | Reputational | scandal, misconduct, controversy, harassment |
| NF | No Adverse Finding | result is neutral or clearly unrelated to the subject |

### Step 3 — Rate Overall Risk

Consider only **Confirmed** and **Probable** findings when assigning the overall risk:

- **High** — at least one Confirmed adverse finding, or multiple Probable findings across different categories.
- **Medium** — one or more Probable adverse findings, or a Confirmed finding that is historical (>10 years ago) or minor in severity.
- **Low** — no adverse findings, or all findings are Excluded.

### Step 4 — Produce the Structured Report

Use exactly the following format:

---

### Subject
- **Name:** {{.Name}}
- **Context:** (fill from context, or "Not provided")
- **Date of Analysis:** (today's date)

### Sources Reviewed
List the query types that returned results (e.g. "fraud, corruption, criminal, processo, escândalo"). If no queries returned results, state that explicitly.

### Adverse Findings

| # | Category | Finding Summary | Source | Date | Relevance |
|---|----------|----------------|--------|------|-----------|
| 1 | | | | | Confirmed / Probable / Excluded |

> If no adverse findings were identified, write: *No adverse findings identified across all queries.*

### Overall Risk Rating
**High / Medium / Low** — one-sentence justification referencing the key finding(s).

### Recommendation
Choose one and provide a brief explanation:
- **Clear** — no further action required.
- **Further investigation required** — specify what additional information would resolve ambiguity.
- **Escalate to compliance team** — describe the confirmed risk requiring escalation.

### Analyst Notes
Note any caveats, search limitations, language barriers, or disambiguation uncertainties. Flag if results were sparse or if the name is very common.

---
