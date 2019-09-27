<script>
    import {link} from 'svelte-spa-router'

    export var users = null;

    function formatDate(timeStamp) {
        var userDate = new Date();
        userDate.setTime(timeStamp * 1000)
        return userDate.toUTCString();
    }
</script>

{#if users}
<table class="table">
    <thead>
        <tr>
            <th>Email</th>
            <th>Created</th>
        </tr>
    </thead>
    <tbody>
        {#each users as user}
        <tr>
            <td><a href="/user/{user.id}" use:link>{user.email}</a></td>
            <td>{formatDate(user.created)}</td>
        </tr>
        {/each}
    </tbody>

</table>
{/if}
