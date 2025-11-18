---
name: Code Review Checklist
about: Checklist for code reviews
title: 'Code Review: '
labels: 'code-review'
assignees: ''

---

## Code Review Checklist

### General
- [ ] The code is well-structured and follows the project's coding standards
- [ ] Variable and function names are clear and descriptive
- [ ] Comments are clear, accurate, and explain why rather than what
- [ ] Duplicated code has been eliminated
- [ ] The code is as simple as possible while still being clear

### Functionality
- [ ] All new functionality is covered by unit tests
- [ ] Error handling is appropriate and consistent
- [ ] Edge cases have been considered and handled appropriately
- [ ] The code handles invalid input gracefully
- [ ] Performance implications have been considered

### Security
- [ ] Input validation is in place for all user-provided data
- [ ] Sensitive data is not logged
- [ ] Proper authentication and authorization checks are in place
- [ ] Dependencies are up to date and secure

### Testing
- [ ] Unit tests cover the new/changed functionality
- [ ] Integration tests are updated if necessary
- [ ] Tests are clear and verify the expected behavior
- [ ] Test data is realistic and covers edge cases

### Documentation
- [ ] Code-level documentation (comments) is sufficient
- [ ] Public APIs are documented
- [ ] README files are updated if necessary
- [ ] Configuration files are documented

### Deployment
- [ ] Database migrations are included if needed
- [ ] Configuration changes are backward compatible
- [ ] The change can be deployed without downtime if required
- [ ] Rollback plan is considered