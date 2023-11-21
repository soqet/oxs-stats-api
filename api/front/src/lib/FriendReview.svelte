<script lang="ts">
  import type { FriendsResponse } from "./schema.ts"

  export let rev: FriendsResponse;

</script>
Игр в одном лобби: {rev?.games_in_same_lobby} 
<div>Среднее место за последние игры: </div>
<div>{rev?.avg_places.first} {rev?.avg_places.second}</div>
<div>Среднее место за последние игры в одном лобби:</div>
<div>{rev?.same_lobby_avg_places.first} {rev?.same_lobby_avg_places.second}</div>
<div>Получено птс за игры в одном лобби</div>
<div>{rev?.pts_gained.first} {rev?.pts_gained.second}</div>