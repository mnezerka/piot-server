<script>
    import {onMount} from 'svelte';
    import {token} from '../stores'
    import Users from '../components/Users.svelte';

    let error = null;
    let users = null;
    let fetching = false;

    onMount(() => {
        fetchUsers();
    })

    function fetchUsers()
    {
        fetching = true;
        error = false;
        users = null;

        var requestBody = {
            query: "{users {email, created}}"
        }

        fetch('http://localhost:9096/query', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + $token,
            },
            body: JSON.stringify(requestBody)
        })
            .then(function(response) {
                if (!response.ok) {
                    error = 'Request failed (' + response.statusText + ')';
                    fetching = false;
                    throw response
                }
                return response.json()
            })
            .then(function(data) {
                fetching = false;
                users = data.data.users
            })
            .catch(function(response) {
                error = 'Request failed (' + response.statusText + ')';
                fetching = false;
            });
    }
</script>


<h1 class="title">Users</h1>

{#if fetching}
    <progress class="progress is-small is-primary" max="100">15%</progress>
{:else}
    {#if error}
        <div class="notification is-danger">
            {error}
        </div>
    {:else}
        <Users users={users}/>
    {/if}
{/if}
