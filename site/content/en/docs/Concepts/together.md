+++
title= "Putting it all together"
description= """How parigot unifies all the concepts in this section."""
date='2023-07-04'
weight= 4
+++

We have discussed in this Concepts section three concepts that are clearly somewhat related:

	1. WASM
	2. Marshaling and unmarshaling
	3. Remote proceedure calls

... but parigot makes them brothers.

parigot lets you do the following things (see if you can spot the correspondence!):

	1. Define an application of many services where the interfaces between the services is specified with the protobuf IDL.
	2. Program that app using the programming language of your choice.
	3. Test and debug your application as single program that has multiple
	services within it. 
	4. Deploy your application as a constellation of services
	that are separated by a network.

If you skipped the previous sections, the connection is 1 in the second
list corresponds to 2 in the first list. Similarly, 2 corresponds to 1, and
3 and 4 of the second use are the two variants of 3 discussed earlier.

Happy microservicing!
