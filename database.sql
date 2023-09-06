-- DROP SCHEMA public;

CREATE SCHEMA public AUTHORIZATION postgres;

-- DROP TYPE public."ProfileType";

CREATE TYPE public."ProfileType" AS ENUM (
	'HODLER',
	'RECYCLER',
	'WASTE_GENERATOR');

-- DROP TYPE public."ResidueType";

CREATE TYPE public."ResidueType" AS ENUM (
	'GLASS',
	'METAL',
	'ORGANIC',
	'PAPER',
	'PLASTIC');
-- public."User" definition

-- Drop table

-- DROP TABLE public."User";

CREATE TABLE public."User" (
	id text NOT NULL,
	"authUserId" text NOT NULL,
	"profileType" public."ProfileType" NOT NULL,
	"lastLoginDate" timestamp(3) NULL,
	"createdAt" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updatedAt" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"phoneNumber" text NOT NULL,
	email text NOT NULL,
	"name" text NOT NULL,
	CONSTRAINT "User_pkey" PRIMARY KEY (id)
);
CREATE UNIQUE INDEX "User_authUserId_key" ON public."User" USING btree ("authUserId");


-- public."_prisma_migrations" definition

-- Drop table

-- DROP TABLE public."_prisma_migrations";

CREATE TABLE public."_prisma_migrations" (
	id varchar(36) NOT NULL,
	checksum varchar(64) NOT NULL,
	finished_at timestamptz NULL,
	migration_name varchar(255) NOT NULL,
	logs text NULL,
	rolled_back_at timestamptz NULL,
	started_at timestamptz NOT NULL DEFAULT now(),
	applied_steps_count int4 NOT NULL DEFAULT 0,
	CONSTRAINT "_prisma_migrations_pkey" PRIMARY KEY (id)
);


-- public."Form" definition

-- Drop table

-- DROP TABLE public."Form";

CREATE TABLE public."Form" (
	"userId" text NOT NULL,
	"createdAt" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updatedAt" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
	id text NOT NULL,
	"isFormAuthorizedByAdmin" bool NULL,
	"walletAddress" text NULL,
	"formMetadataUrl" text NULL,
	CONSTRAINT "Form_pkey" PRIMARY KEY (id),
	CONSTRAINT "Form_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON DELETE CASCADE ON UPDATE CASCADE
);


-- public."Issue" definition

-- Drop table

-- DROP TABLE public."Issue";

CREATE TABLE public."Issue" (
	id varchar(40) NOT NULL,
	userid text NOT NULL,
	total_issuance numeric(65, 30) NOT NULL,
	report jsonb NOT NULL,
	wallet varchar(80) NOT NULL,
	createdat timestamp NOT NULL,
	CONSTRAINT "Issue_pkey" PRIMARY KEY (id),
	CONSTRAINT issue_userid_fkey FOREIGN KEY (userid) REFERENCES public."User"(id) ON DELETE CASCADE ON UPDATE CASCADE
);


-- public."Allocation" definition

-- Drop table

-- DROP TABLE public."Allocation";

CREATE TABLE public."Allocation" (
	txhash varchar(80) NOT NULL,
	"percent" numeric(65, 30) NOT NULL,
	wallet varchar(80) NOT NULL,
	issueid varchar(40) NOT NULL,
	total numeric(65, 30) NOT NULL,
	CONSTRAINT allocation_fk FOREIGN KEY (issueid) REFERENCES public."Issue"(id)
);


-- public."Certificate" definition

-- Drop table

-- DROP TABLE public."Certificate";

CREATE TABLE public."Certificate" (
	id varchar(40) NOT NULL,
	txhash varchar(80) NOT NULL,
	certid int8 NOT NULL,
	createdat timestamptz NOT NULL,
	formid text NOT NULL,
	walletperform varchar(80) NOT NULL,
	towallet varchar(80) NOT NULL,
	CONSTRAINT certificate_pk PRIMARY KEY (id),
	CONSTRAINT form_fk FOREIGN KEY (formid) REFERENCES public."Form"(id)
);
CREATE UNIQUE INDEX certificate_formid_idx ON public."Certificate" USING btree (formid);


-- public."Document" definition

-- Drop table

-- DROP TABLE public."Document";

CREATE TABLE public."Document" (
	id text NOT NULL,
	"formId" text NOT NULL,
	"videoFileName" text NULL,
	"invoicesFileName" _text NULL,
	"createdAt" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updatedAt" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"residueType" public."ResidueType" NOT NULL,
	amount numeric(65, 30) NOT NULL,
	CONSTRAINT "Document_pkey" PRIMARY KEY (id),
	CONSTRAINT "Document_formId_fkey" FOREIGN KEY ("formId") REFERENCES public."Form"(id) ON DELETE CASCADE ON UPDATE CASCADE
);


-- public.issue_certificate definition

-- Drop table

-- DROP TABLE public.issue_certificate;

CREATE TABLE public.issue_certificate (
	issueid varchar(40) NOT NULL,
	certificateid varchar(40) NOT NULL,
	issued bool NOT NULL,
	CONSTRAINT issue_certificate_fk FOREIGN KEY (certificateid) REFERENCES public."Certificate"(id),
	CONSTRAINT issue_certificate_fk_1 FOREIGN KEY (issueid) REFERENCES public."Issue"(id)
);
