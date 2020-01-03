<script>
    import {link} from 'svelte-spa-router'
    import {formatDate} from '../utils'

    export var things = null;
</script>

{#if things}
<table class="table">
    <thead>
        <tr>
            <th>Name</th>
            <th>Alias</th>
            <th>Type</th>
            <th>Created</th>
            <th>Enabled</th>
            <th>Organization</th>
        </tr>
    </thead>
    <tbody>
        {#each things as thing}
        <tr>
            <td><a href="/thing/{thing.id}" use:link>{thing.name}</a></td>
            <td>{thing.alias}</td>
            <td>{thing.type}</td>
            <td>{formatDate(thing.created)}</td>
            <td>
            {#if thing.enabled}
                <span class="icon has-text-success"><i class="fas fa-check-square"></i></span>
            {:else}
                <span class="icon has-text-danger"><i class="fas fa-ban"></i></span>
            {/if}
            </td>
            <td>{#if thing.org}{thing.org.name}{/if}</td>
        </tr>
        {/each}
    </tbody>

</table>
{/if}
