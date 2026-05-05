# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

## [1.0.0] - 2026-05-04

### Added

- TCP group chat server accepting up to 10 concurrent clients
- ASCII art welcome banner and username prompt on connect
- Username validation with re-prompt on empty input
- Message history replayed to new joiners before join notification
- Join and leave broadcast notifications to all other clients
- Timestamped message format `[YYYY-MM-DD HH:MM:SS][name]:message`
- Empty and whitespace-only message filtering
- Thread-safe client registry with atomic check-and-add to enforce the 10-client limit
- Default port `8989` with optional command-line override
- Port validation (digits only, range 1–65535)
- Usage message on invalid arguments
- Full test suite with race detector support
