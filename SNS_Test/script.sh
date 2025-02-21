aws ce put-budget \
  --budget-name "FreeTierLimit" \
  --time-unit MONTHLY \
  --budget-type COST \
  --limit-amount 0.01 \
  --limit-unit USD

