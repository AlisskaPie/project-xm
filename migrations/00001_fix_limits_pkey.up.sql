ALTER TABLE IF EXISTS public.company
    ALTER COLUMN id SET NOT NULL,
    ALTER COLUMN name TYPE character varying(15),
    ALTER COLUMN description TYPE character varying(3000),
    ALTER COLUMN description DROP NOT NULL,
    ALTER COLUMN type SET NOT NULL,
    ADD PRIMARY KEY (id);
