<!-- The title of your incident. Make sure the title is a h1 (single #)-->
# Love Lost Globally: Jerry Develops Malicious App

<!-- The date which the incident started on. The no letters after the numbers please 1, 2, 3 NOT 1st, 2nd, 3rd -->
## Date: July 1, 2020

<!-- The owner of the post mortem, responsible for following up on actions -->
## Owner: Morty Smith

A short description of the event. Rick help develop the malicious app of an innocent alien.
Hostile aliens almost take over the planets water supply.

## Timeline

| Time | Event |
| --- | --- |
| 9:16 | Breakfast. Rick introduces alien. "Do not develop my app" is tattooed on forehead |
| 10:37 | Jerry begins app development with alien |
| 12:30 | App released |
| 12:34 | Morty realises aliens app is released |
| 15:36 | Lovefinderrz reaches 100,000 users |
| 18:44 | Jerry and Morty install paywall, number of users rapidly decreases |
| 20:03 | No app users remain |

<!-- Crucial metrics to agree on. Format: x unit[, x smaller_unit] -->
<!-- Units can be seconds, minutes, hours, days. Use full unit (1 second, not 1s) -->
<!-- Severity can be on your own scale, it is tracked as a category rather than a metric -->
<!-- One example: 1 = Service down completely, 2 = Service hindered for many users, 3 = Service hindered for some -->
## Metrics

| Metric | Time |
| --- | --- |
| Severity | 1 |
| Time To Detect | 4 minutes |
| Time To Resolve | 6 hours, 14 minutes |
| Total Downtime | 6 hours, 28 minutes | <!-- Detect + Resolve, Sanity check. -->

## Cause of the Problem

Alien with malicious intent invited into the house. Family members not informed of the severity
of alien app. Forehead tattoo documentation inadequate.


## Corrective Actions with Owners

* All house members must be debriefed before high risk aliens are brought into house [Rick]
  - Enforced using debrief document created in family-process repo
* App review process to require 3 reviewers before release [Jerry]
