# Progression Discord Bot

This project aims to simplify our MTG progression league, by enabling players to easily manage their card pools using a simple discord bot.
Players can join the progression league.
The admin(s) can launch the league. This will generate pairings for the first round. 
Players can self-report the match results. 
Once all matches have been reported, every player is awarded a wild card. Every losing player is awarded a wild pack. Refer to the commands section for details on how to redeem cards and packs.
Finally, the admin starts the next round: the next set becomes available, every player is given 10 wild packs, and new pairings are generated.

At any point, players can get their current card pool and wild card/pack count.

## Commands
<details>
<summary>Players</summary>

These commands are available to all players or those looking to join the league.

<details>
<summary>
<code>/join</code> - Join the upcoming league
</summary>

Join the league. Generates a card pool of 10 packs for the joining player.

**Syntax:**
`/join`

**Arguments:**
None

**Restriction:**

The command will fail if:
- a league is ongoing
</details>
<details>
<summary>
<code>/drop</code> - Drop from the league
</summary>
Removes the player from the current league.
If the player is part of a match, which hasn't been reported on yet, it will be set to 2:0 for the opponent.

**Syntax:**
`/drop`

**Arguments:**
None

**Restriction:**

The command will fail if:
- the given user does not play in the current league.
</details>

<details>
<summary>
<code>/report</code> - Report match results
</summary>

The games won always refer to the reporting player. 
The opponent does not need to report the same match.

**Syntax:**
`/report <games_won> <games_lost> <draws>`

**Arguments:**
- `<games_won>`is the number of games in the match won by the reporting player.
- `<games_lost>`is the number of games in the match won by the opponent of the reporting player.
- `<draws>`is the number of games in the match ending in a draw.

**Restriction:**

The command will fail if:
- no league is ongoing
- the match result has already been reported
</details>

<details>
<summary>
<code>/redeem</code> - Redeem wild cards & packs
</summary>

The `/redeem` command has two different sub commands for redeeming individual cards or packs.

<details>
<summary>
<code>/redeem card</code> - Turning wild cards into specific cards
</summary>

Spend a wild card to get add a specific card from the unlocked sets to your card pool.

**Syntax:**
`/redeem card <set_code> <collector_number>`

**Arguments:**
- `<set_code>`is a valid set code of an already unlocked set in the current league.
- `<collector_number>`is the collector number of the card in the given set you want to add to your card pool.

**Restriction:**

The command will fail if:
- no league is ongoing
- the user doesn't play in the current league
- the user has no more wild cards
- the given set code is not valid or not unlocked (yet) for the current league
- the given collector number is not valid or does not exist in the given set
</details>

<details>
<summary>
<code>/redeem pack</code> - Turning wild packs into packs of random cards
</summary>

Spend one or more wild packs to add that many packs of random cards from one of the unlocked sets to your card pool.

**Syntax:**
`/redeem pack <set_code> <count>`

**Arguments:**
- `<set_code>`is a valid set code of an already unlocked set in the current league.
- `<count>`is the number of packs of the given set you want to add to your card pool.

**Restriction:**

The command will fail if:
- no league is ongoing
- the user doesn't play in the current league
- the user has insufficient wild packs
- the given set code is not valid or not unlocked (yet) for the current league
</details>
</details>
<details>
<summary>
<code>/pool</code> - Get a list of all cards in your card pool
</summary>

Get a list of all cards in your personal card pool.

**Syntax:**
`/pool`

**Arguments:**
None

**Restriction:**

The command will fail if:
- no league is ongoing
- the user doesn't play in the current league
</details>
<details>
<summary>
<code>/balance</code> - Check the number of wild cards & packs available to you
</summary>

Get the number of wild cards & packs you can still redeem.

**Syntax:**
`/balance`

**Arguments:**
None

**Restriction:**

The command will fail if:
- no league is ongoing
- the user doesn't play in the current league
</details>
</details>

<details>
<summary>Admins</summary>

These commands are only available to administrators of the league.

<details>
<summary>
<code>/start</code> - Start a new league
</summary>
Start the league with all players that have joined so far.
The given set will be the first available set to redeem wild packs and cards for.

**Syntax:**
`/start <set_code>`

**Arguments:**
- `<set_code>`is a valid MTG set code of the first set to make available to all players.

**Restriction:**

The command will fail if:
- a league is ongoing
- an invalid set code is given
</details>

<details>
<summary>
<code>/next</code> - Start the next round
</summary>
Start the next round in the current league.
The given set will become available to redeem wild packs and cards for.

**Syntax:**
`/next <set_code>`

**Arguments:**
- `<set_code>`is a valid MTG set code of the next set to make available to all players.

**Restriction:**

The command will fail if:
- no league is active
- an invalid set code is given
- at least one match hasn't been reported on yet
</details>

<details>
<summary>
<code>/force_drop</code> - Remove a player from the league
</summary>
Remove a player from the current league.
If the player is part of a match, which hasn't been reported on yet, it will be set to 2:0 for the opponent.

**Syntax:**
`/force_drop <username>`

**Arguments:**
- `<username>`is a valid discord handle of a player in the current league.

**Restriction:**

The command will fail if:
- an invalid username is given
- the given user does not play in the current league.
</details>

<details>
<summary>
<code>/force_report</code> - Set a match result
</summary>

Report on a match from the perspective of the given player.

**Syntax:**
`/force_report <username> <games_won> <games_lost> <draws>`

**Arguments:**
- `<username>`is a valid discord handle of a player in the current league.
- `<games_won>`is the number of games the given player has won.
- `<games_lost>`is the number of games the given player has lost.
- `<draws>`is the number draws in the match.

**Restriction:**

The command will fail if:
- an invalid username is given
- the given user does not play in the current league.
- the given user was not part of a match (e.g. in case of an uneven number of players).
</details>

<details>
<summary>
<code>/ban</code> - Ban a card from the current league
</summary>

Add a card to the ban list.

**Syntax:**
`/ban <cardname>`

**Arguments:**
- `<cardname>`is a valid MTG card name.

**Restriction:**

The command will fail if:
- no league is active
- the given cardname does not match a valid card
- the given cardname is not part of any of the unlocked sets
- the given cardname is already on the ban list
</details>
<details>
<summary>
<code>/unban</code> - Unban a card from the current league
</summary>

Remove a card from the ban list.

**Syntax:**
`/unban <cardname>`

**Arguments:**
- `<cardname>`is a valid MTG card name.

**Restriction:**

The command will fail if:
- no league is active
- the given cardname is not on the ban list
</details>
</details>

## Open Questions
- Matches are a best of three?
- How are scores calculated?