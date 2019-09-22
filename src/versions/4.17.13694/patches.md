Due to the large amount of changes in this firmware, many patches are not needed
anymore, and others need rewriting. I will also be cleaning up the less used
patches (if anyone needs them, PM me or open an issue). Some new patches will
also be added.

I am also doing a cleanup of the nickel patches to make them more readable and
future-proof. Warning: this includes patch renames.

This list applies to 4.16.13337, 4.17.13651, and 4.17.13694 (they will be released
together due to the fast pace of the recent updates).

In the patch files, I will be marking patches which need testing (TODO17:TEST)
and rewriting (TODO17:REWRITE). These markers will all need to be dealt with and
removed before the final release.

Quite a few CSS patches (not listed here) also needed additional values for the
Libra H2O.

## Added
- geek1011 / Remove title from reading header/footer
- jackie_w / Reduce new header/footer height

### TODO

### Done

---

## Waiting (for usual update)
- GeoffR / My 24 line spacing values
- GeoffR / ePub fixed/adjustable top/bottom margins

---

## Rewritten

### TODO
- GeoffR / Custom font sizes - needs to be updated for Storm/Libra, offsets have changed

**Unplanned (i.e. not a release blocker, nobody has taken these yet):**
- jcn363 / Changing the info panel in full size screensaver - will need to be rewritten if people are still interested in this patch.
- GeoffR / Keyboard patches - the keyboard has completely changed, @geek1011 is not familiar with these patches and doesn't have time to update them.

### Done
- oren64 / Increase size of kepub chapter progress chart - needed to be rewritten
- geek1011 / Both buttons go next - code has been refactored, but seems it might be simpler
- oren64 / Increase the view details container size - will need to be rewritten if people are still interested in this patch - jackie_w rewrote this as "Increase book details synopsis area"
- GeoffR / Custom footer (page number text) - completely different in new reader view, will need a complete rewrite if someone still wants it - jackie_w rewrote this as "Custom header/footer page number text"
- jackie_w / Dictionary text font-family/font-size/line-height - beta - jackie_w is making some changes and taking it out of beta
- oren64 / Increase the cover size in library - jackie_w is rewrote and simplified this to use the new kobopatch instructions
- oren64 / Custom font to collection and author titles - will need to be rewritten if people are still interested in this patch - jackie_w rewrote this patch as "Custom font for Collection/Author header titles"

---

## Removed

### Not needed
- GeoffR / Clock display duration - only applies to old reader view, new one has a persistent clock in the status bar.
- GeoffR / Fix three KePub fullScreenReading bugs - the bugs seem to have been fixed since 4.11.11911.
- GeoffR / Always display chapter name on navigation menu - new reader view does this by default.
- GeoffR / Fix reading stats/author name cut off when series is showing - now fixed in the firmware.
- geek1011 / Remove extra space on selection menu - doesn't make a difference anymore.
- geek1011 / Set reading footer height - only applies to old reader.
- jackie_w / Custom menubar - reduce height by 33% - doesn't do what it used to, not needed for the reader view anymore (jackie_w *might* make another)
- jackie_w / Custom menubar - reduce height by 50% - doesn't do what it used to, not needed for the reader view anymore (jackie_w *might* make another)
- GeoffR / Custom reading footer style - doesn't apply to new reader view (would need to be rewritten), and it's not really needed anymore
- GeoffR / Disable reading footer - no effect on new reader view.

### Cleaned up
- geek1011 / Rename settings - I made it for something I wanted to test, no need for it anymore, I don't think anyone uses it. 

---

## Renamed
- geek1011 / Enable rotation on all devices -> DeveloperSettings - ForceAllowLandscape
- geek1011 / Set slide to unlock -> PowerSettings - UnlockEnabled
- jcn363 / Wikipedia Search language -> Change Wikipedia search language
- jcn363 / Dictionary, remove long winded closest match notification -> Shorten dictionary entry not found message
- GeoffR / New home screen increasing cover size -> Increase home screen cover size
