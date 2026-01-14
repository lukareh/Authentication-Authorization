Task: Build a tiny in-memory API simulation protecting resources
Implement a simple backend-like program (no actual server needed):

Create an array of documents like:
[

  {id: 1, owner: "userA", content: "A's secret"},

  {id: 2, owner: "userB", content: "B's secret"}

]

 Ask the user to input:
Their username
A document ID they want to access

Program should:
Fetch the document
Check if document.owner === username
If yes: print content
If no: print "Access Denied"

Bonus:
Add a second user role (admin) who can access all documents.

Implementation:
This project implements a document access control system using Go as the backend with an HTTP REST API and a simple HTML frontend. 
The backend maintains an in-memory array of documents and validates user access based on ownership, with special admin privileges. 
Users can request document access through a web interface, and the system checks if the user owns the document or has admin role before granting access. 
The application demonstrates basic authorization concepts with role-based access control (RBAC), where regular users can only access their own documents while admins have unrestricted access to all resources.