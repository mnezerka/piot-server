<script>
    import {formatDate} from '../utils';
    import {link} from 'svelte-spa-router'

    export var customers = null;
</script>

{#if customers}
<table class="table">
    <thead>
        <tr>
            <th>Email</th>
            <th>Created</th>
        </tr>
    </thead>
    <tbody>
        {#each customers as customer}
        <tr>
            <td><a href="/customer-view/{customer.id}" use:link>{customer.name}</a></td>
            <td>{formatDate(customer.created)}</td>
        </tr>
        {/each}
    </tbody>

</table>
{/if}
