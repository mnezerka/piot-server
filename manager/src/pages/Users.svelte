<script>
    //import TestCasesList from '../components/TestCasesList.svelte';

    let error = null;
    let users = null;
    let fetching = false;

    /*
    let field_types = [
        {id: 'twbc_guid', name: 'ID'},
        {id: 'twbc_title', name: 'Title'}
    ]

    function onSearch()
    {
        fetching = true;
        error = false;
        test_cases = null;

        let query_url = 'https://gtp.wdf.sap.corp/sap/bc/pi-tools/qte/query_twbc?format=json&' +
            selected_field.id + '=' + search;


        fetch(query_url)
            .then(function(response) {
                if (!response.ok) {
                    error = 'Request failed (' + response.statusText + ')';
                    throw response
                }
                fetching = false;
                return response.json()
            })
            .then(function(data) {
                if (data.DATA) {
                    data = data.DATA;
                }
                test_cases = data
            })
            .catch(function(response) {
                error = 'Request failed (' + response.statusText + ')';
                fetching = false;
            });
    }
    */
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
        users
    {/if}
{/if}
