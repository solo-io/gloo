import {
  DataSource,
  ObjectRef,
  Selector
} from 'proto/dev-portal/api/dev-portal/v1/common_pb';
import { Portal } from 'proto/dev-portal/api/grpc/admin/portal_pb';
import { ObjectMeta, Time } from 'proto/dev-portal/api/grpc/common/common_pb';
import {
  PortalSpec,
  StaticPage,
  KeyScope,
  PortalStatus,
  CustomStyling,
  KeyScopeStatus
} from 'proto/dev-portal/api/dev-portal/v1/portal_pb';

export function createDataSourceClassFromObject(
  dataSourceObj: DataSource.AsObject | undefined
): DataSource | undefined {
  if (!dataSourceObj) {
    return undefined;
  } else {
    let dataSourceClass = new DataSource();

    dataSourceClass.setFetchUrl(dataSourceObj?.fetchUrl);

    if (!!dataSourceObj?.configMap) {
      let configMap = new DataSource.ConfigMapData();
      configMap.setName(dataSourceObj.configMap.name);
      configMap.setNamespace(dataSourceObj.configMap.namespace);
      configMap.setKey(dataSourceObj.configMap.key);
      dataSourceClass.setConfigMap(configMap);
    }

    dataSourceClass.setInlineBytes(dataSourceObj?.inlineBytes);

    dataSourceClass.setInlineString(dataSourceObj?.inlineString);

    return dataSourceClass;
  }
}

function createSelectorClassFromObject(
  selectorObj: Selector.AsObject | undefined
): Selector | undefined {
  let selectorClass = new Selector();

  if (!selectorObj) {
    return undefined;
  }

  selectorObj.matchLabelsMap.forEach(([key, val]) => {
    selectorClass.getMatchLabelsMap().set(key, val);
  });

  return selectorClass;
}
function createObjectRefClassFromObject(
  objectRefObj: ObjectRef.AsObject
): ObjectRef {
  let objRefClass = new ObjectRef();

  objRefClass.setNamespace(objectRefObj.namespace);
  objRefClass.setName(objectRefObj.name);

  return objRefClass;
}

export function createPortalClassFromObject(
  portalObject: Portal.AsObject
): Portal {
  let portalAsClass = new Portal();

  // METADATA
  if (!!portalObject.metadata) {
    let objectMeta = new ObjectMeta();

    if (!!portalObject.metadata?.creationTimestamp) {
      let timestamp = new Time();
      timestamp.setNanos(portalObject.metadata?.creationTimestamp.nanos);
      timestamp.setSeconds(portalObject.metadata?.creationTimestamp.seconds);
      objectMeta.setCreationTimestamp(timestamp);
    }
    objectMeta.setName(portalObject.metadata?.name || 'BAD NAME');
    objectMeta.setNamespace(
      portalObject.metadata?.namespace || 'BAD NAMESPACE'
    );
    objectMeta.setResourceVersion(
      portalObject.metadata?.resourceVersion || 'BAD VERSION'
    );

    portalAsClass.setMetadata(objectMeta);
  }

  // PORTAL SPEC
  if (!!portalObject.spec) {
    let portalSpec = new PortalSpec();

    portalSpec.setDescription(portalObject.spec.description || '');
    portalSpec.setDisplayName(portalObject.spec.displayName || '');
    portalSpec.setDomainsList(portalObject.spec.domainsList || []);
    portalSpec.setBanner(
      createDataSourceClassFromObject(portalObject.spec.banner)
    );
    portalSpec.setFavicon(
      createDataSourceClassFromObject(portalObject.spec.favicon)
    );
    portalSpec.setPrimaryLogo(
      createDataSourceClassFromObject(portalObject.spec.primaryLogo)
    );
    portalSpec.setPublishApiDocs(
      createSelectorClassFromObject(portalObject.spec.publishApiDocs)
    );

    portalSpec.setStaticPagesList(
      portalObject.spec.staticPagesList.map(pageObj => {
        let pageClass = new StaticPage();
        pageClass.setContent(createDataSourceClassFromObject(pageObj.content));
        pageClass.setDescription(pageObj.description);
        pageClass.setName(pageObj.name);
        pageClass.setNavigationLinkName(pageObj.navigationLinkName);
        pageClass.setPath(pageObj.path);
        return pageClass;
      })
    );

    portalSpec.setKeyScopesList(
      portalObject.spec.keyScopesList.map(keyScopeObj => {
        let keyScopeClass = new KeyScope();
        keyScopeClass.setApiDocs(
          createSelectorClassFromObject(keyScopeObj.apiDocs)
        );
        keyScopeClass.setDescription(keyScopeObj.description);
        keyScopeClass.setName(keyScopeObj.name);
        keyScopeClass.setNamespace(keyScopeObj.namespace);
        return keyScopeClass;
      })
    );

    if (!!portalObject.spec.customStyling) {
      let customStyling = new CustomStyling();
      customStyling.setBackgroundColor(
        portalObject.spec.customStyling.backgroundColor
      );
      customStyling.setButtonColorOverride(
        portalObject.spec.customStyling.buttonColorOverride
      );
      customStyling.setDefaultTextColor(
        portalObject.spec.customStyling.defaultTextColor
      );
      customStyling.setNavigationLinksColorOverride(
        portalObject.spec.customStyling.navigationLinksColorOverride
      );
      customStyling.setPrimaryColor(
        portalObject.spec.customStyling.primaryColor
      );
      customStyling.setSecondaryColor(
        portalObject.spec.customStyling.secondaryColor
      );
      portalSpec.setCustomStyling(customStyling);
    }

    portalAsClass.setSpec(portalSpec);
  }

  // PORTAL STATUS
  if (!!portalObject.status) {
    let portalStatus = new PortalStatus();

    portalStatus.setObservedGeneration(portalObject.status.observedGeneration);
    portalStatus.setReason(portalObject.status.reason);
    portalStatus.setPublishUrl(portalObject.status.publishUrl);
    portalStatus.setState(portalObject.status.state);

    portalStatus.setApiDocsList(
      portalObject.status.apiDocsList.map(apiDocObj =>
        createObjectRefClassFromObject(apiDocObj)
      )
    );

    portalStatus.setKeyScopesList(
      portalObject.status.keyScopesList.map(keyScopeStatusObj => {
        let keyScopeStatus = new KeyScopeStatus();

        keyScopeStatus.setName(keyScopeStatusObj.name);
        keyScopeStatus.setAccessibleApiDocsList(
          keyScopeStatusObj.accessibleApiDocsList.map(accApiDoc =>
            createObjectRefClassFromObject(accApiDoc)
          )
        );
        keyScopeStatus.setProvisionedKeysList(
          keyScopeStatusObj.provisionedKeysList.map(provKey =>
            createObjectRefClassFromObject(provKey)
          )
        );

        return keyScopeStatus;
      })
    );
    /*
    keyScopesList: Array<KeyScopeStatus.AsObject>,*/

    portalAsClass.setStatus(portalStatus);
  }

  return portalAsClass;
}
