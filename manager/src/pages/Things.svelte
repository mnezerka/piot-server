<script>
    import {onMount} from 'svelte';
    import {authenticated, token} from '../stores'
    import {push} from 'svelte-spa-router'
    import {gql} from '../utils.js';
    import ThingList from '../components/ThingList.svelte';

    let error = null;
    let things = null;
    let fetching = false;

    onMount(() => {
        if (!$authenticated) { push("/login"); }
        fetchThings();

    })

    async function fetchThings()
    {
        fetching = true;
        error = false;

        try {
            let data = await gql({query: "{things {id, name, alias, type, enabled, created, org {name}}}"});
            things = data.things;
        } catch(error) {
            error = 'Request failed (' + error + ')';
        }

        fetching = false;
    }

</script>

<h1 class="title">Things</h1>

{#if fetching}
    <progress class="progress is-small is-primary" max="100">15%</progress>
{:else}
    {#if error}
        <div class="notification is-danger">
            {error}
        </div>
    {:else}
        <ThingList things={things}/>
    {/if}
{/if}
