API Endpoints
-------------



Get Seed
----------

.. http:get:: /keys/seed
    Results:

   .. sourcecode:: js
    your  seed 

Create Key
............

.. http:post:: /keys
    :string name: **Required**. the name of a user
    :string password: **Required**. the password of a user
    :string seed q: **Required**. The search query
    Results:

   .. sourcecode:: js
    your address 
 
