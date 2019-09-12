Due to the large amount of changes in this firmware, many patches are not needed
anymore, and others need rewriting. I will also be cleaning up the less used
patches (if anyone needs them, PM me or open an issue). Some new patches will
also be added.

## Added

---

## Rewritten

### Waiting
- GeoffR / Keyboard patches - the keyboard has completely changed, @geek1011 is not familiar with these patches and doesn't have time to update them.
- GeoffR / Disable reading footer - changed in the new reader view.
- GeoffR / Reduce top/bottom page spacer - might need to be rewritten

### Done

---

## Removed

### Not needed
- GeoffR / Clock display duration - only applies to old reader view, new one has a persistent clock in the status bar.
- GeoffR / Fix three KePub fullScreenReading bugs - the bugs seem to have been fixed since 4.11.11911.
- GeoffR / Always display chapter name on navigation menu - new reader view does this by default.
- GeoffR / Fix reading stats/author name cut off when series is showing - now fixed in the firmware.

### Not possible
- GeoffR / Custom footer (page number text) - completely different in new reader view, will need a complete rewrite if someone still wants it.
- GeoffR / Custom reading footer style - see note in nickel.yaml.

### Cleaned up
- geek1011 / Rename settings - I made it for something I wanted to test, no need for it anymore, I don't think anyone uses it. 

---

## Renamed
- geek1011 / Enable rotation on all devices -> DeveloperSettings - ForceAllowLandscape
- geek1011 / Set slide to unlock -> PowerSettings - UnlockEnabled
