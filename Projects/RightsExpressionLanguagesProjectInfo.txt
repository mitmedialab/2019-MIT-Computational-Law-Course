Exercise: Represent MMA Rights in REL

1. As RightsML style data structure
2. Extra credit: in ODRL JSON or RDF syntax
3. Extra extra credit: represent all rights for interactive music DSPs as ODRL profile
a. Stream, conditional download, permanent download
b. Composition and sound recording rights

* http://dev.iptc.org/RightsML-Simple-Example-Geographic
* https://www.w3.org/community/odrl/implementations/ 
* https://www.w3.org/TR/2018/REC-odrl-vocab-20180215/
* https://www.w3.org/TR/2018/REC-odrl-model-20180215/ 
* https://github.com/nitmws/odrl-wprofile-evaltest1/blob/master/evaluator/evaluator.js



The License as RightsML/ODRL Policy XML document


<o:Policy uid="http://example.com/RightsML/policy/idGeog1"
 type="http://www.w3.org/ns/odrl/2/Set"
 profile="https://iptc.org/std/RightsML/odrl-profile/"
 xmlns:o="http://www.w3.org/ns/odrl/2/">
   <o:permission>
      <o:asset uid="urn:newsml:example.com:20120101:180106-999-000013"
      
      /for sound recording use the ISRC and for composition use ISWC 
      / issue - you are given the sound recording then you need to find the matching composition and HOPE the ISRC still exists... but good enough for a classroom exercise 
       relation="http://www.w3.org/ns/odrl/2/target"/>
      <o:action name="http://www.w3.org/ns/odrl/2/distribute"/>
      / Use DDEX codes (slide 29)
      / Codes for rights granted (distribution channels), e.g.:
      / OnDemandStream
      / Conditional-Download
      / PermanentDownload
      
      / DDEX: MM-0891 - Release Notification 2PD3.7 (1).pdf (PAGE 116 use
     /rights)
     / 
https://kb.ddex.net/download/attachments/2294037/MM-0891%20-%20Release%20Notification%202PD3.7.pdf?version=1&modificationDate=1378118185872
       
      / Note: p. 74 lists DistributionChannelType codes, which are essentially
      / rights to be conferred by the licensee on end users. E.g.
      /OnDemandStream, ConditionalDownload.

      <o:constraint name="http://www.w3.org/ns/odrl/2/spatial"
       operator="http://www.w3.org/ns/odrl/2/eq"
       rightOperand="http://cvx.iptc.org/iso3166-1a3/DEU"/>
      <o:party uid="http://example.com/cv/party/epa"
       function="http://www.w3.org/ns/odrl/2/assigner"/>
      <o:party uid="http://example.com/cv/partygroup/epapartners"
       function="http://www.w3.org/ns/odrl/2/assignee"
       type="http://www.w3.org/ns/odrl/2/PartyCollection" />
   </o:permission>
</o:Policy>
