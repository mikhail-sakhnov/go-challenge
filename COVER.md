Cover
=====

Timeouts
________
First of all, timeout for 500ms. As far as I know, go lang program running on standard linux kernel can't guarantee hard real-time timeouts.
So, I used standard context.Context structure with setted deadline, but I can't protected the whole code with this timeout, at least I would have to write answer to ResponseWriter with results I'd have at the moment of timeout expiration.
Because of it, I'd come to the next scheme:
    - I have only one handler in the service.
    - Business related logic totally protected by context Timeout, by default 450 ms.
    - Remained 50ms are for encoding json and writing answer. I can decrease this value, by serializing each intermediate result, but I afraid it would cost a much and it wouldn't give much.

Urls
____
Service got urls via get parameters. There is no any mentions of escaping this parameters, so possible, url with get parameter that would passed unurlencoded would break url parsing.
But assuming that examples has no encoding and no urls with parameters I didn't implement urlencoding.
I can download urls serial, but it is bad decision - one slow url will destroy the whole results set.
For concurrent approach to download urls there is no any mentions about worker limit. Having no limits is bad, so I decided to limit workers by MAX_CPU or passed argument.

Also, there is no mentions about behaviour in case of double urls, so I just skip doubles.

Merging and storing results
_______________
There are few possible approaches to store results:
- I can store all retrieved numbers as slice and resort slice each time after receiving new part of numbers (it would cost N * nlogn for where N is for urls count and n is already got numbers count)
- I can store them all in unsorted slice and sort only before writing answer (but, it would increase 50ms value which is not context deadline protected)
- In both approaches I should also prevent doubles, which is possible by having some kind of hash (e.g. map[int]struct{}) for O(1) checks if we have such number.
- More optimal approach is to use sort.Search function which allows to find place in slice where I should insert new element, using binary search. Of course, it means that I should have to maintain sorted order in slice.
- I can use some self sorting data structure.

I chose the last approach with avl tree. AVL tree is good for two reasons:
- After building tree I can perform dfs traverse on it and go sorted slice of all tree values
- It's automatically prevent doubles on insert.

Unfortunatelly, there is no AVL tree implementation in golang library, so I implemented it by myself.

Sync
____
At the top, each goroutine is protected from leak by select with context.Done channel.
context.Done channel automatically closed by http library after request handling or on deadline, if specified.
Each fetcher also guarded by wait group.

There are next entities:
    - handler is for serving http requests.
    - Service is structure with business related logic (downloading and merging).
    - Url sanitizer is for skipping invalid urls and prevent doubles.
    - Fetcher func is for fetching remote numbers.
    - AVL tree is for storing results.
    - Response object is for encoding and decoding json.

On each incoming request, handler creates instances of url sanitizer and service. Service starts downloading pipeline, providing two channels for listeting on to the handler:
    - Done channel. This channel would be closed after processing each url.
    - Tick channel. This channel get messages with []int after each new chunk of numbers completed.
Handler use select for listening on this channels and also ctx.Done for handling timeout.
On each Tick message, handler stores the received []int in response object.
If ctx.Done or service.Done happened, handler encodes response object and write response to the RW.

