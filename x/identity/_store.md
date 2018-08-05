# Stores 

## Owners 
- Prefix Key Space: OwnersKey 
- Key/Sort: Ident Address Then Owner Address
- Value: empty

## Owner Count
- Prefix Key Space: OwnerCountPrefix 
- Key/Sort: Ident Address
- Value: Number

## Certs 
- Prefix Key Space: CertsKey
- Key/Sort: Ident Address Then Property Name Then Issuer Address
- Value: Cert Object

## Trusts 
- Prefix Key Space: TrustsKey
- Key/Sort: Validator Address Then Ident Address
