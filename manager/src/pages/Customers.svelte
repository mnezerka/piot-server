<script>
    import {onMount} from 'svelte';
    import {token} from '../stores'
    import Customers from '../components/Customers.svelte';
    import {gql} from '../utils.js';
    import {push} from 'svelte-spa-router'

    let error = null;
    let customers = null;
    let fetching = false;

    onMount(() => {
        fetchCustomers();
    })

    async function fetchCustomers() {
        fetching = true;
        error = false;
        customers = null;

        try {
            let data = await gql({query: "{customers {name, created}}"});
            customers = data.customers;
        } catch(error) {
            error = 'Request failed (' + error + ')';
        }

        fetching = false;
    }

    function onAdd() {
        push('/customer-add');
    }

</script>


<h1 class="title">Customers</h1>

{#if fetching}
    <progress class="progress is-small is-primary" max="100">15%</progress>
{:else}
    {#if error}
        <div class="notification is-danger">
            {error}
        </div>
    {:else}
        <Customers customers={customers}/>
        <button class="button" on:click={onAdd}>Add</button>
    {/if}
{/if}
