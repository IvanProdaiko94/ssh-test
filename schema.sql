create table games
(
  /* x-o-ox-xx */
  board char(9) not null,
  /* UUIDv4 */
  id varchar(36) not null,
  /* in case of adding some long status */
  status varchar(20) not null
);
