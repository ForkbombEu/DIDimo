# 🏗️ Building blocks

The following system diagrams provide a visual overview of the DIDimo platform, illustrating its architecture and key components. These diagrams are based on the C4 model, which is a widely-used approach for visualizing software architecture.

While we follow the principles of the C4 model, our implementation is somewhat loose, focusing on clarity and relevance to the specific context of DIDimo. The diagrams capture different levels of abstraction, from the overall system context down to detailed component interactions within the platform. This structured approach helps in understanding how various parts of the system interact and contribute to its overall functionality.


## System Context diagrams

Shows the interaction between external actors (developers, service providers, governmental bodies, etc.) and the DIDimo system.

@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Container.puml
!define I https://raw.githubusercontent.com/tupadr3/plantuml-icon-font-sprites/master/govicons
!include I/users.puml
!include I/user_politician.puml
!include I/user_suit.puml
!include I/ribbon.puml
!include I/presenter.puml

LAYOUT_WITH_LEGEND()

title System Context Diagram for DIDimo

Person(dev, "Developer", "Submits credential issuers and runs compliance checks.")
Person(sp, "Service Provider", "Manages and publishes compliance results.")
Person(gov, "Government IT Manager", "Ensures compliance with standards.", $sprite="user_suit")
Person(eu, "EU Official", "Generates compliance reports.", $sprite="user_politician")
Person(cto, "CTO", "Browses and compares credential services.")
Person(enduser, "End User", "Browses verified credential issuers.", $sprite="users")
Person(researcher, "Researcher", "Analyzes compliance data.", $sprite="presenter")
Person(sb, "Standardization Body", "Evaluates standards alignment.", $sprite="ribbon")

System(didimo, "DIDimo", "Platform for verifying compliance of decentralized identity services.")

Rel(dev, didimo, "Submits and checks compliance")
Rel_U(sp, didimo, "Manages and publishes results")
Rel_U(gov, didimo, "Generates compliance reports")
Rel_R(eu, didimo, "Generates compliance reports")
Rel_R(cto, didimo, "Browses and compares services")
Rel_L(enduser, didimo, "Browses verified issuers")
Rel_L(researcher, didimo, "Analyzes compliance data")
Rel(sb, didimo, "Evaluates standards alignment")
Lay_U(enduser, researcher)
Lay_U(cto, eu)
@enduml

## Container diagrams
Illustrates the main containers within the DIDimo system (API Gateway, Compliance Engine, Dashboard, etc.) and their interactions.

@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Container.puml
!define I https://raw.githubusercontent.com/tupadr3/plantuml-icon-font-sprites/master/govicons
!define FA6 https://raw.githubusercontent.com/tupadr3/plantuml-icon-font-sprites/master/font-awesome-4
!include I/user_politician.puml
!include I/user_suit.puml
!include I/ribbon.puml
!include I/presenter.puml

LAYOUT_WITH_LEGEND()

title Container Diagram for DIDimo

Person(dev, "Developer")
Person(sp, "Service Provider")
Person(gov, "Government IT Manager", $sprite="user_suit")
Person(eu, "EU Official", $sprite="user_politician")
Person(cto, "CTO")
Person(enduser, "End User")
Person(researcher, "Researcher", $sprite="presenter")
Person(sb, "Standardization Body", $sprite="ribbon")

System_Boundary(didimo, "DIDimo") {
    Container(api_gateway, "API Gateway", "golang", "Handles API requests and CLI submissions.")
    Container(compliance_engine, "Compliance Engine", "Zenroom / Slangroom", "Performs compliance checks and runs periodic tasks.")
    Container(database, "Database", "MongoDB", "Stores compliance data, user data, and reports.")
    Container(queue_manager, "Queue Manager", "RabbitMQ / golang", "Manages long-running compliance tasks.")
    Container(reporting_service, "Reporting Service", "golang", "Generates and exports compliance reports.")
    Container(external_services_api, "External Services API", "APIs", "Connects to external compliance check services.")
    Container(dashboard, "Dashboard", "TypeScript/Svelte", "User interface for managing services and viewing results.")
    Container(comparison_tool, "Marketplace/Comparison Tool", "TypeScript/Svelte", "Tool for comparing credential services.")
}

Rel(dev, api_gateway, "Submits credential issuers/checks via API/CLI")
Rel(sp, dashboard, "Manages and publishes results")
Rel(enduser, comparison_tool, "Browses verified issuers")
Rel(cto, dashboard, "Browses and compares services")

Rel_U(sb, reporting_service, "Evaluates standards alignment")
Rel_U(gov, reporting_service, "Generates compliance reports")
Rel_U(eu, reporting_service, "Generates compliance reports")
Rel_U(researcher, reporting_service, "Accesses compliance data")

Rel(api_gateway, compliance_engine, "Processes compliance checks")
Rel_L(compliance_engine, database, "Stores compliance results")
Rel(queue_manager, compliance_engine, "Manages long-running tasks")
Rel(compliance_engine, external_services_api, "Uses external services for checks")
Rel(dashboard, compliance_engine, "Schedules periodic checks")
Rel(dashboard, database, "Manages user data")
Rel(comparison_tool, database, "Accesses compliance data")
Rel_U(reporting_service, database, "Fetches data for reports")

@enduml



## Component Diagram
Focuses on the internal components of the Compliance Engine, showing how the various modules work together.

@startuml
!include <C4/C4_Component>

!define DEVICONS https://raw.githubusercontent.com/tupadr3/plantuml-icon-font-sprites/master/devicons
!define FONTAWESOME https://raw.githubusercontent.com/tupadr3/plantuml-icon-font-sprites/master/font-awesome-5
!include DEVICONS/go.puml
!include FONTAWESOME/users.puml

title Component Diagram for Compliance Engine

Container_Boundary(didimo_backend, "DIDimo Backend - Compliance Engine") {
    Component(compliance_engine, "Compliance Engine", "Slangroom/ncr", "Performs all compliance checks.")
    Component(standards_checker, "Standards Checker", "Step CI/slangroom", "Validates against various identity standards.")
    Component(debugging_tool, "Debugging Tool", "Slangroom/golang", "Helps developers resolve compliance issues.")
    Component(periodic_checker, "Periodic Checker", "golang", "Schedules and runs periodic checks.")
    Component(external_service_integration, "External Service Integration", "REST", "Connects to external compliance services via APIs.")
}

Rel(compliance_engine, standards_checker, "Uses for compliance validation")
Rel(compliance_engine, debugging_tool, "Uses for debugging and issue resolution")
Rel(compliance_engine, periodic_checker, "Uses for scheduling periodic checks")
Rel(compliance_engine, external_service_integration, "Uses for external compliance checks")

@enduml


