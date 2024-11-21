
0.6.2 / 2024-11-21
==================

  * feat: handle graceful shutdown (#64)
  * fix: expose cleanup errors (#62)
  * feat: add option to ignore missing deployment (#33)
  * fix(cleanup): inactivate all deployment related to PR environment
  * chore: lots of nix and go module updates

0.6.1 / 2022-06-28
==================

  * ci: bump go versions
  * add dependabot to automate upgrades

0.6.0 / 2022-06-28
==================

  * fix: list existing deployments in the same environment (#14)
  * feat: switch to flake, add fmt, lint (#13)
  * feat(cleanup): destroy github deployments related to PR (#12)
  * feat(cleanup): filter deployments based on a given PR list (#11)
  * feat(cleanup): handle script arguments (#10)
  * fix: remove duplicate arguments command (#9)
  * feat: run from subdirectories in git repo (#8)
  * feat(please): also display stdout (#6)
  * feat(please): handle scripts arguments (#7)
  * Merge pull request #4 from jfroche/feat/provide-environment-url
  * ci: build/test using different go versions
  * feat(please): set environment URL using new cli flag

0.5.0 / 2022-01-15
==================

  * bump all the dependencies
  * add support for GitHub Actions

0.4.n / 2020-12-18
==================

  * release with goreleaser
  * bye Travis-CI
  * Create go.yml
  * update dependencies
  * please: trim the environmentURL
  * README: update with mdsh
  * restore the original environment name logic
  * forward API failures to the client
  * convert to go modules
  * gitsrc: wrap errors for better reporting
  * fix CLI parsing
  * Use branch as deployment ref by default (#1)
  * add missing CHANGELOG.md
  * fix: handle pull-request URLs

0.3.0 / 2019-04-29
==================

  * wrap the github token in zimbatm/go-secretvalue

0.2.0 / 2019-04-29
==================

  * add support for Circle CI

0.1.0 / 2018-06-03
==================

  * cleanup: fix the opened PR list
  * report to the right commit ID in Travis
  * fixes bug with Travis CI
  * relax the deployment names
  * fix github slug parsing
  * Hi

