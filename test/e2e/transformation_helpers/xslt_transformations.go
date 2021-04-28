package transformation_helpers

const JsonToXmlTransform = `
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform" xmlns="http://www.w3.org/2005/xpath-functions" version="3.0">
    <xsl:output indent="yes" />
    <xsl:template match="/">
        <C_MKTSEGMENT_ROLLUP xmlns="">
            <xsl:apply-templates select="json-to-xml(.)/*"/>
        </C_MKTSEGMENT_ROLLUP>
    </xsl:template>
    <xsl:template match="map" xpath-default-namespace="http://www.w3.org/2005/xpath-functions">
        <C_MKTSEGMENT_SUMMARY xmlns="">
            <C_MKTSEGMENT><xsl:value-of select="string[@key='market_segment']"/></C_MKTSEGMENT>
            <C_ACCTBAL_TOTAL><xsl:value-of select="number[@key='account_balance_total']"/></C_ACCTBAL_TOTAL>
            <C_ACCTBAL_TREND><xsl:value-of select="string[@key='account_balance_trend']"/></C_ACCTBAL_TREND>
            <C_NATIONKEYS>
                <xsl:for-each select="array[@key='nation_ids']/number">
                    <C_NATIONKEY><xsl:value-of select="."/></C_NATIONKEY>
                </xsl:for-each>
            </C_NATIONKEYS>
        </C_MKTSEGMENT_SUMMARY>
    </xsl:template>
</xsl:stylesheet>
`

const XmlToJsonTransform = `<?xml version="1.0" encoding="UTF-8"?>
          <xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
            version="3.0">
            <xsl:output method="text"/>
            <xsl:template match="/">
                <xsl:value-of select="xml-to-json(., map { 'indent' : true() })"/>
            </xsl:template>
          </xsl:stylesheet>`

const XmlToXmlTransform = `<?xml version="1.0" encoding="UTF-8"?>
        <xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform" 
            xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:max="http://www.ibm.com/maximo" exclude-result-prefixes="soapenv max">
            <xsl:output method="xml" indent="yes"/>
            <xsl:strip-space elements="*" />

            <xsl:template match="max:PO">
                <SyncX1POMATREC>
                    <X1POMATRECSet>
                        <PO>
                            <SITEID><xsl:value-of select="max:SITEID" /></SITEID>
                            <PONUM><xsl:value-of select="max:PONUM" /></PONUM>
                            <REVISIONNUM><xsl:value-of select="max:REVISIONNUM" /></REVISIONNUM>
                            <POID><xsl:value-of select="POID" /></POID>
                            <RECEIPT>
                                <POLINENUM><xsl:value-of select="POLINENUM" /></POLINENUM>
                                <ITEMNUM><xsl:value-of select="ITEMNUM" /></ITEMNUM>
                            </RECEIPT>
                        </PO>
                    </X1POMATRECSet>
                </SyncX1POMATREC>
            </xsl:template>
        </xsl:stylesheet>`

const AutomobileJson = `[
    {
        "market_segment": "AUTOMOBILE",
        "account_balance_total": 75577498.12,
        "account_balance_trend": "DECLINING",
        "nation_ids": [1, 2]
    },
    {
        "market_segment": "MACHINERY",
        "account_balance_total": 5544522866.83,
        "account_balance_trend": "INCREASING",
        "nation_ids": [4, 2]
    }
]`

const AutomobileXml = `<?xml version="1.0" encoding="UTF-8"?>
<C_MKTSEGMENT_ROLLUP>
   <C_MKTSEGMENT_SUMMARY>
      <C_MKTSEGMENT>AUTOMOBILE</C_MKTSEGMENT>
      <C_ACCTBAL_TOTAL>75577498.12</C_ACCTBAL_TOTAL>
      <C_ACCTBAL_TREND>DECLINING</C_ACCTBAL_TREND>
      <C_NATIONKEYS>
         <C_NATIONKEY>1</C_NATIONKEY>
         <C_NATIONKEY>2</C_NATIONKEY>
      </C_NATIONKEYS>
   </C_MKTSEGMENT_SUMMARY>
   <C_MKTSEGMENT_SUMMARY>
      <C_MKTSEGMENT>MACHINERY</C_MKTSEGMENT>
      <C_ACCTBAL_TOTAL>5544522866.83</C_ACCTBAL_TOTAL>
      <C_ACCTBAL_TREND>INCREASING</C_ACCTBAL_TREND>
      <C_NATIONKEYS>
         <C_NATIONKEY>4</C_NATIONKEY>
         <C_NATIONKEY>2</C_NATIONKEY>
      </C_NATIONKEYS>
   </C_MKTSEGMENT_SUMMARY>
</C_MKTSEGMENT_ROLLUP>
`

const CarsJson = `
  { "🚗" : 
    [ 
      { "doors" : "🚗",
        "price" : "6Ü" },
      
      { "doors" : "5ü",
        "price" : "13L" } ] }`

const CarsXml = `<?xml version="1.0" encoding="UTF-8"?>
  <map xmlns="http://www.w3.org/2005/xpath-functions">
      <array key="🚗">
          <map>
              <string key="doors">🚗</string>
              <string key="price">6Ü</string>
          </map>
          <map>
              <string key="doors">5ü</string>
              <string key="price">13L</string>
          </map>
      </array>
  </map>`

const SoapMessageXml = `<?xml version="1.0" encoding="utf-8" ?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:max="http://www.ibm.com/maximo">
   <soapenv:Header/>
   <soapenv:Body>
      <max:SyncX1POMATREC>
         <max:X1POMATRECSet>
            <max:PO action="AddChange">
               <max:SITEID>BEDFORD</max:SITEID>
               <max:PONUM>TEST42</max:PONUM>
               <max:REVISIONNUM>0</max:REVISIONNUM>
            </max:PO>
         </max:X1POMATRECSet>
      </max:SyncX1POMATREC>
   </soapenv:Body>
</soapenv:Envelope>`

const TransformedSoapXml = `<?xml version="1.0" encoding="UTF-8"?>
<SyncX1POMATREC>
   <X1POMATRECSet>
      <PO>
         <SITEID>BEDFORD</SITEID>
         <PONUM>TEST42</PONUM>
         <REVISIONNUM>0</REVISIONNUM>
         <POID/>
         <RECEIPT>
            <POLINENUM/>
            <ITEMNUM/>
         </RECEIPT>
      </PO>
   </X1POMATRECSet>
</SyncX1POMATREC>
`
