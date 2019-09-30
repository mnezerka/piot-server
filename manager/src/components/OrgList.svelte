<script>
    import {formatDate} from '../utils';
    import {link} from 'svelte-spa-router'

    export var orgs = null;
</script>

{#if orgs}
<table class="table">
    <thead>
        <tr>
            <th>Email</th>
            <th>Created</th>
        </tr>
    </thead>
    <tbody>
        {#each orgs as org}
        <tr>
            <td><a href="/org/{org.id}" use:link>{org.name}</a></td>
            <td>{formatDate(org.created)}</td>
        </tr>
        {/each}
    </tbody>

</table>
{/if}
