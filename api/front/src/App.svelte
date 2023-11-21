<script lang="ts">
import FriendReview from "./lib/FriendReview.svelte";
import type { FriendsResponse, ResponseSchema } from "./lib/schema";



let first_id: number | undefined = undefined;
let second_id: number | undefined = undefined;
let isRevVisible = false;
let frev: Promise<ResponseSchema>;

async function getFriends(first: number | undefined, second: number | undefined) {
    if (first == undefined || second == undefined) {
        return
    }
    isRevVisible = true
    let resp = await fetch("http://localhost:9999/api/review/friends", {
        method: "POST",
        body: JSON.stringify({
            first_id: first,
            second_id: second,
        }),
    })
    frev = resp.json()
}


</script>

<label for="first_id"> First id</label>
<input type="number" id="first_id" bind:value={first_id}>
<label for="second_id"> Second id</label>
<input type="number" id="second_id" bind:value={second_id}>
<button on:click={() => getFriends(first_id, second_id)}> Проверить силу их дружбы </button>


{#if isRevVisible}
<div>
    {#await frev}
        загрузка
    {:then rev} 
        <FriendReview rev={rev?.data} />
    {:catch}
        <p style="color: red"> анлак </p>
    {/await}
</div>
{/if}