<script>
    import {formatDate} from '../utils';
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
            <td>{customer.name}</td>
            <td>{formatDate(customer.created)}</td>
        </tr>
        {/each}
    </tbody>

</table>
{/if}
