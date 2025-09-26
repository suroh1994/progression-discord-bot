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
<summary>
<code>/join</code>
</summary>

Join the league. Generates a card pool of 10 packs for the joining player.

**Syntax:**
`/join`

**Arguments:**
None

**Permission:**
Everyone

**Restriction:**
Only available until the league starts.
</details>

<details>
<summary>
<code>/start</code>
</summary>
</details>

<details>
<summary>
<code>/report</code>
</summary>

The games won always refer to the reporting player. The opponent does not need to report the same match.

**Syntax:**
`/report <games won> <games lost> <draws>`

**Arguments:**
- `<games won>`is the number of games in the match won by the reporting player.
- `<games lost>`is the number of games in the match won by the opponent of the reporting player.
- `<draws>`is the number of games in the match ending in a draw.

**Permission:**
Players

**Restriction:**
Only available during a round. Only one player per pairing can report the result. Any further reports will be ignored.
</details>

<details>
<summary>
<code>/redeem</code>
</summary>

The `/redeem` command has two different sub commands for redeeming individual cards or packs.

<details>
<summary>
<code>/redeem card</code>
</summary>
</details>

<details>
<summary>
<code>/redeem pack</code>
</summary>
</details>
</details>

<details>
<summary>
<code>/cardpool</code>
</summary>
</details>

<details>
<summary>
<code>/balance</code>
</summary>
</details>

<details>
<summary>
<code>/next</code>
</summary>
</details>