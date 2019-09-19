<script>
    import {token, authenticated} from '../stores.js'
    import {push} from 'svelte-spa-router'
    import {gql} from '../utils.js';
    import {onMount} from 'svelte';

    export var params;

    let error = null;
    let fetching = false;
    let customer = null;

    onMount(async () => {
        if (params.id) {
            await fetchCustomer(params.id);
        } else {
            error = 'No customer specified';
        }
    })

    async function fetchCustomer(name) {
        fetching = true;
        error = false;
        customer = null;

        try {
            let data = await gql({query: `{customer (name: "${name}") {name, description, created}}`});
            customer = data.customer;
        } catch(error) {
            error = 'Request failed (' + error + ')';
        }

        fetching = false;
    }

    function onEdit() {
        push(`/customer-edit/${params.id}`)
    }

</script>

<h1 class="title">Customer {params.id}</h1>

{#if fetching}
    <progress class="progress is-small is-primary" max="100">15%</progress>
{:else}
    {#if error}
        <div class="notification is-danger">
            {error}
        </div>
    {:else}
        {#if customer}
            <div class="content">
                <ul>
                    <li>Name: {customer.name}</li>
                    <li>Description: {customer.description}</li>
                    <li>Created: {customer.created}</li>
                </ul>
            </div>
            <button class="button" on:click={onEdit}>Edit</button>
        {/if}
    {/if}
{/if}


