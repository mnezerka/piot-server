<script>
    export let error = null;
</script>

{#if error}
    <div class="notification is-danger">

    {#if error instanceof Array}
        {#each error as e}
            {e.message}<br/>
        {/each}
    {:else}
        {error}
    {/if}
    </div>
{/if}
